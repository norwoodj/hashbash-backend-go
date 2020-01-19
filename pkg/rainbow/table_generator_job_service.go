package rainbow

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"sync/atomic"

	"github.com/norwoodj/hashbash-backend-go/pkg/dao"
	"github.com/norwoodj/hashbash-backend-go/pkg/metrics"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type TableGenerateJobConfig struct {
	ChainBatchSize int64
	NumThreads     int
}

type TableGeneratorJobService struct {
	jobConfig              TableGenerateJobConfig
	rainbowTableService    dao.RainbowTableService
	rainbowChainService    dao.RainbowChainService
	chainGenerationSummary *prometheus.SummaryVec
	chainWriteSummary      *prometheus.SummaryVec
	chainsCreatedCounter   *prometheus.CounterVec
}

func NewRainbowTableGeneratorJobService(
	jobConfig TableGenerateJobConfig,
	rainbowChainService dao.RainbowChainService,
	rainbowTableService dao.RainbowTableService,
	chainGenerationSummary *prometheus.SummaryVec,
	chainWriteSummary *prometheus.SummaryVec,
	chainsCreatedCounter *prometheus.CounterVec,
) *TableGeneratorJobService {
	return &TableGeneratorJobService{
		jobConfig:              jobConfig,
		rainbowChainService:    rainbowChainService,
		rainbowTableService:    rainbowTableService,
		chainGenerationSummary: chainGenerationSummary,
		chainWriteSummary:      chainWriteSummary,
		chainsCreatedCounter:   chainsCreatedCounter,
	}
}

func (service *TableGeneratorJobService) saveRainbowChains(
	rainbowTable *model.RainbowTable,
	chainList []model.RainbowChain,
) error {
	labeledSummary := service.chainWriteSummary.
		With(metrics.GetRainbowTableMetricLabels(rainbowTable, len(chainList)))

	timer := prometheus.NewTimer(labeledSummary)
	defer timer.ObserveDuration()

	return service.rainbowChainService.CreateRainbowChains(rainbowTable.ID, chainList)
}

func (service *TableGeneratorJobService) runChainGeneratorThread(
	rainbowTable *model.RainbowTable,
	chainGeneratorService *chainGeneratorService,
	batchesRemaining *int64,
	errorChannel chan error,
) {
	defer close(errorChannel)
	log.Debugf("Spawned thread to generate chains for rainbow table %d, in batches of %d", rainbowTable.ID, service.jobConfig.ChainBatchSize)
	chainList := make([]model.RainbowChain, service.jobConfig.ChainBatchSize)
	chainLength := int(rainbowTable.ChainLength)

	for atomic.AddInt64(batchesRemaining, -1) >= 0 {
		timer := prometheus.NewTimer(service.chainGenerationSummary.
			With(metrics.GetRainbowTableMetricLabels(rainbowTable, int(service.jobConfig.ChainBatchSize))))

		for i := 0; i < int(service.jobConfig.ChainBatchSize); i++ {
			startPlaintext := chainGeneratorService.NewRandomString(rainbowTable.CharacterSet, rainbowTable.PasswordLength)
			chainList[i] = chainGeneratorService.generateRainbowChain(startPlaintext, chainLength)
		}

		timer.ObserveDuration()

		err := service.saveRainbowChains(rainbowTable, chainList)
		if err != nil {
			errorChannel <- err
		}

		err = service.rainbowTableService.IncrementRainbowTableChainsGenerated(rainbowTable.ID, service.jobConfig.ChainBatchSize)
		if err != nil {
			errorChannel <- err
		}

		service.chainsCreatedCounter.
			With(metrics.GetRainbowTableMetricLabels(rainbowTable, int(service.jobConfig.ChainBatchSize))).
			Add(float64(service.jobConfig.ChainBatchSize))
	}

	log.Debugf("No batches remaining to be generated for rainbow table %d, exiting", rainbowTable.ID)
}

func (service *TableGeneratorJobService) checkAndUpdateRainbowTableStatus(rainbowTable *model.RainbowTable) error {
	if rainbowTable.Status != model.StatusQueued {
		return fmt.Errorf("cannot generate a rainbow table that is not in the %s state", model.StatusQueued)
	}

	err := service.rainbowTableService.UpdateRainbowTableStatusAndGenerateStarted(rainbowTable.ID, model.StatusStarted)
	if err != nil {
		return fmt.Errorf("failed to update rainbow table status to %s: %s", model.StatusStarted, err)
	}

	return nil
}

func (service *TableGeneratorJobService) calculateNumBatches(rainbowTable *model.RainbowTable) int64 {
	batchesRemaining := rainbowTable.NumChains / service.jobConfig.ChainBatchSize

	if rainbowTable.NumChains%service.jobConfig.ChainBatchSize > 0 {
		batchesRemaining += 1
	}

	return batchesRemaining
}

func (service *TableGeneratorJobService) spawnChainGeneratorThreads(
	rainbowTable *model.RainbowTable,
	hashFunctionProvider HashFunctionProvider,
	errorChannels []chan error,
) {
	batchesRemaining := service.calculateNumBatches(rainbowTable)

	for i := 0; i < service.jobConfig.NumThreads; i++ {
		chainGeneratorService := newChainGeneratorService(
			hashFunctionProvider.NewHashFunction(),
			getDefaultReductionFunctionFamily(int(rainbowTable.PasswordLength), rainbowTable.CharacterSet),
			i+1,
		)

		go service.runChainGeneratorThread(
			rainbowTable,
			chainGeneratorService,
			&batchesRemaining,
			errorChannels[i],
		)
	}
}

func (service *TableGeneratorJobService) awaitChainGenerationCompletionOrErrors(
	rainbowTable *model.RainbowTable,
	errorChannels []chan error,
) error {
	threadsInError := 0
	for _, e := range errorChannels {
		err := <-e
		if err != nil {
			threadsInError += 1
		}
	}

	if threadsInError == service.jobConfig.NumThreads {
		err := service.rainbowTableService.UpdateRainbowTableStatus(rainbowTable.ID, model.StatusFailed)
		if err != nil {
			return fmt.Errorf(
				"generation failed for rainbow table %d and failed to update rainbow table status to %s: %s",
				rainbowTable.ID,
				model.StatusFailed,
				err,
			)
		}

		return fmt.Errorf("all generate threads failed, failing generate job for rainbow table %d", rainbowTable.ID)
	}

	return nil
}

func (service *TableGeneratorJobService) updateFinalChainCountAndStatus(rainbowTable *model.RainbowTable) error {
	finalChainCount, err := service.rainbowChainService.CountChainsForRainbowTable(rainbowTable.ID)
	if err != nil {
		return fmt.Errorf(
			"failed to retrieve final chain count for rainbow table %d: %s",
			rainbowTable.ID,
			err,
		)
	}

	err = service.rainbowTableService.UpdateRainbowTableStatusAndFinalChainCount(
		rainbowTable.ID,
		model.StatusCompleted,
		finalChainCount,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to update rainbow table status to %s for rainbow table %d: %s",
			model.StatusCompleted,
			rainbowTable.ID,
			err,
		)
	}

	return nil
}

func (service *TableGeneratorJobService) RunGenerateJobForTable(rainbowTableId int16) error {
	rainbowTable, err := service.rainbowTableService.FindRainbowTableById(rainbowTableId)

	if gorm.IsRecordNotFoundError(err) {
		return fmt.Errorf("rainbow table with ID %d not found, cannot generate", rainbowTableId)
	} else if err != nil {
		return err
	}

	err = service.checkAndUpdateRainbowTableStatus(&rainbowTable)
	if err != nil {
		return err
	}

	hashFunctionProvider, err := GetHashFunctionProvider(rainbowTable.HashFunction)
	if err != nil {
		return err
	}

	errorChannels := initializeErrorChannels(service.jobConfig.NumThreads)
	service.spawnChainGeneratorThreads(&rainbowTable, hashFunctionProvider, errorChannels)

	err = service.awaitChainGenerationCompletionOrErrors(&rainbowTable, errorChannels)
	if err != nil {
		return err
	}

	return service.updateFinalChainCountAndStatus(&rainbowTable)
}
