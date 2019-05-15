package rainbow

import (
	"fmt"
	"math/rand"
	"strings"
	"sync/atomic"

	"github.com/norwoodj/hashbash-backend-go/pkg/model"
	"github.com/norwoodj/hashbash-backend-go/pkg/service"
	log "github.com/sirupsen/logrus"
)

type TableGenerateJobConfig struct {
	ChainBatchSize int64
	NumThreads     int
}

type TableGeneratorJobService struct {
	jobConfig           TableGenerateJobConfig
	rainbowTableService service.RainbowTableService
	rainbowChainService service.RainbowChainService
}

func NewTableGeneratorJobService(
	rainbowChainService service.RainbowChainService,
	rainbowTableService service.RainbowTableService,
	jobConfig TableGenerateJobConfig,
) *TableGeneratorJobService {
	return &TableGeneratorJobService{
		jobConfig:           jobConfig,
		rainbowChainService: rainbowChainService,
		rainbowTableService: rainbowTableService,
	}
}

func randomString(characterSet *string, stringLength int64) string {
	stringBuilder := strings.Builder{}
	characterSetSize := len(*characterSet)
	var i int64
	for i = 0; i < stringLength; i++ {
		index := rand.Intn(characterSetSize)
		stringBuilder.WriteByte((*characterSet)[index])
	}

	return stringBuilder.String()
}

func (service *TableGeneratorJobService) runChainGeneratorThread(
	rainbowTable model.RainbowTable,
	chainGeneratorService *chainGeneratorService,
	batchesRemaining *int64,
	errorChannel chan error,
) {
	defer close(errorChannel)
	log.Debugf("Spawned thread to generate chains for rainbow table %d, in batches of %d", rainbowTable.ID, service.jobConfig.ChainBatchSize)
	chainList := make([]model.RainbowChain, service.jobConfig.ChainBatchSize)
	chainLength := int(rainbowTable.ChainLength)

	for atomic.AddInt64(batchesRemaining, -1) >= 0 {
		for i := 0; i < int(service.jobConfig.ChainBatchSize); i++ {
			startPlaintext := randomString(&rainbowTable.CharacterSet, rainbowTable.PasswordLength)
			chainList[i] = chainGeneratorService.generateRainbowChain(startPlaintext, chainLength)
		}

		err := service.rainbowChainService.CreateRainbowChains(rainbowTable.ID, chainList)
		if err != nil {
			errorChannel <- err
		}

		err = service.rainbowTableService.IncrementRainbowTableChainsGenerated(rainbowTable.ID, service.jobConfig.ChainBatchSize)
		if err != nil {
			errorChannel <- err
		}
	}

	log.Debugf("No batches remaining to be generated for rainbow table %d, exiting", rainbowTable.ID)
}

func (service *TableGeneratorJobService) checkAndUpdateRainbowTableStatus(rainbowTable model.RainbowTable) error {
	if rainbowTable.Status != model.StatusQueued {
		return fmt.Errorf("cannot generate a rainbow table that is not in the %s state", model.StatusQueued)
	}

	err := service.rainbowTableService.UpdateRainbowTableStatus(rainbowTable.ID, model.StatusStarted)
	if err != nil {
		return fmt.Errorf("failed to update rainbow table status to %s: %s", model.StatusStarted, err)
	}

	return nil
}

func (service *TableGeneratorJobService) calculateNumBatches(rainbowTable model.RainbowTable) int64 {
	batchesRemaining := rainbowTable.NumChains / service.jobConfig.ChainBatchSize

	if rainbowTable.NumChains % service.jobConfig.ChainBatchSize > 0 {
		batchesRemaining += 1
	}

	return batchesRemaining
}

func (service *TableGeneratorJobService) initializeErrorChannels() []chan error {
	errorChannels := make([]chan error, service.jobConfig.NumThreads)

	for i := range errorChannels {
		errorChannels[i] = make(chan error)
	}

	return errorChannels
}

func (service *TableGeneratorJobService) spawnChainGeneratorThreads(
	rainbowTable model.RainbowTable,
	batchesRemaining *int64,
	hashFunctionProvider hashFunctionProvider,
	errorChannels []chan error,
) {
	for i := 0; i < service.jobConfig.NumThreads; i++ {
		chainGeneratorService := newChainGeneratorService(
			hashFunctionProvider.newHashFunction(),
			getDefaultReductionFunctionFamily(int(rainbowTable.PasswordLength), rainbowTable.CharacterSet),
		)

		go service.runChainGeneratorThread(
			rainbowTable,
			chainGeneratorService,
			batchesRemaining,
			errorChannels[i],
		)
	}
}

func (service *TableGeneratorJobService) awaitChainGenerationCompletionOrErrors(
	rainbowTable model.RainbowTable,
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

func (service *TableGeneratorJobService) updateFinalChainCountAndStatus(rainbowTable model.RainbowTable) error {
	finalChainCount := service.rainbowChainService.CountChainsForRainbowTable(rainbowTable.ID)
	err := service.rainbowTableService.UpdateRainbowTableFinalChainCount(rainbowTable.ID, finalChainCount)

	if err != nil {
		return fmt.Errorf("failed to update final chain count for rainbow table %d: %s", rainbowTable.ID, err)
	}

	err = service.rainbowTableService.UpdateRainbowTableStatus(rainbowTable.ID, model.StatusCompleted)
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

func (service *TableGeneratorJobService) RunGenerateJobForTable(rainbowTable model.RainbowTable) error {
	err := service.checkAndUpdateRainbowTableStatus(rainbowTable)
	if err != nil {
		return err
	}

	hashFunctionProvider, found := hashFunctionProvidersByName[rainbowTable.HashFunction]
	if !found {
		return fmt.Errorf("invalid hash function specified by rainbow table %s", rainbowTable.HashFunction)
	}

	batchesRemaining := service.calculateNumBatches(rainbowTable)
	errorChannels := service.initializeErrorChannels()

	service.spawnChainGeneratorThreads(rainbowTable, &batchesRemaining, hashFunctionProvider, errorChannels)
	err = service.awaitChainGenerationCompletionOrErrors(rainbowTable, errorChannels)
	if err != nil {
		return err
	}

	return service.updateFinalChainCountAndStatus(rainbowTable)
}
