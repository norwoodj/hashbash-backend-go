package rainbow

import (
	"encoding/hex"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/meirf/gopart"
	"github.com/norwoodj/hashbash-backend-go/pkg/dao"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
)

type TableSearchJobConfig struct {
	SearchHashBatchSize int
	NumThreads          int
}

type TableSearchJobService struct {
	jobConfig                 TableSearchJobConfig
	rainbowChainService       dao.RainbowChainService
	rainbowTableService       dao.RainbowTableService
	rainbowTableSearchService dao.RainbowTableSearchService
}

func NewRainbowTableSearchJobService(
	jobConfig TableSearchJobConfig,
	rainbowChainService dao.RainbowChainService,
	rainbowTableService dao.RainbowTableService,
	rainbowTableSearchService dao.RainbowTableSearchService,
) *TableSearchJobService {
	return &TableSearchJobService{
		jobConfig:                 jobConfig,
		rainbowChainService:       rainbowChainService,
		rainbowTableService:       rainbowTableService,
		rainbowTableSearchService: rainbowTableSearchService,
	}
}

func (service *TableSearchJobService) runChainGenerationThread(
	searchHash []byte,
	rainbowTable model.RainbowTable,
	indexByEndHash map[string]int64,
	hashFunctionProvider HashFunctionProvider,
	currentPossibleIndex *int64,
	indexByEndHashMutex *sync.Mutex,
	waitGroup *sync.WaitGroup,
) {
	defer waitGroup.Done()
	rainbowChainGeneratorService := newChainGeneratorService(
		hashFunctionProvider.NewHashFunction(),
		getDefaultReductionFunctionFamily(int(rainbowTable.PasswordLength), rainbowTable.CharacterSet),
	)

	chainIndex := atomic.AddInt64(currentPossibleIndex, 1)
	for chainIndex < rainbowTable.ChainLength {
		rainbowChain := rainbowChainGeneratorService.generateRainbowChainLinkFromHash(
			searchHash,
			int(chainIndex),
			int(rainbowTable.ChainLength-chainIndex-1),
		)

		hashStringRep := fmt.Sprintf("%x", rainbowChain.hashedPlaintext)

		indexByEndHashMutex.Lock()
		indexByEndHash[hashStringRep] = chainIndex
		indexByEndHashMutex.Unlock()

		chainIndex = atomic.AddInt64(currentPossibleIndex, 1)
	}
}

func (service *TableSearchJobService) getSearchBatches(indexByEndHash map[string]int64) [][]string {
	endHashList := make([]string, 0)
	for endHash := range indexByEndHash {
		endHashList = append(endHashList, endHash)
	}

	endHashSearchBatches := make([][]string, 0)
	for indexRange := range gopart.Partition(len(endHashList), service.jobConfig.SearchHashBatchSize) {
		endHashSearchBatches = append(endHashSearchBatches, endHashList[indexRange.Low:indexRange.High])
	}

	return endHashSearchBatches
}

func (service *TableSearchJobService) runSearchThread(
	rainbowTableId int16,
	searchBatch []string,
	foundChannel chan model.RainbowChain,
) {
	defer close(foundChannel)
	rainbowChain := service.rainbowChainService.FindChainByTableIdAndEndHashIn(rainbowTableId, searchBatch)

	if rainbowChain.EndHash != "" {
		foundChannel <- rainbowChain
	}
}

func (service *TableSearchJobService) generatePlaintextFromFoundEndHash(
	searchHash string,
	foundChannels []chan model.RainbowChain,
	indexByEndHash map[string]int64,
	rainbowChainGeneratorService *chainGeneratorService,
) string {
	for _, foundChannel := range foundChannels {
		rainbowChain := <-foundChannel
		if rainbowChain.EndHash == "" {
			continue
		}

		chainIndex, _ := indexByEndHash[rainbowChain.EndHash]
		plaintextLink := rainbowChainGeneratorService.generateRainbowChainLinkFromPlaintext(
			rainbowChain.StartPlaintext,
			0,
			int(chainIndex+1),
		)

		if fmt.Sprintf("%x", plaintextLink.hashedPlaintext) == searchHash {
			return plaintextLink.plaintext
		}
	}

	return ""
}

func (service *TableSearchJobService) spawnChainGenerationThreads(
	rainbowTableSearch model.RainbowTableSearch,
	rainbowTable model.RainbowTable,
	indexByEndHash map[string]int64,
	errorChannels []chan error,
	hashFunctionProvider HashFunctionProvider,
) error {
	var currentPossibleIndex int64 = -1
	var indexByEndHashMutex sync.Mutex

	searchHash, err := hex.DecodeString(rainbowTableSearch.Hash)

	if err != nil {
		return fmt.Errorf("failed to decode hash string for search: %s", err)
	}

	waitGroup := sync.WaitGroup{}
	for i := 0; i < service.jobConfig.NumThreads; i++ {
		waitGroup.Add(1)
		go service.runChainGenerationThread(
			searchHash,
			rainbowTable,
			indexByEndHash,
			hashFunctionProvider,
			&currentPossibleIndex,
			&indexByEndHashMutex,
			&waitGroup,
		)
	}

	waitGroup.Wait()
	return nil
}

func (service *TableSearchJobService) RunSearchJob(searchId int64) error {
	rainbowTableSearch := service.rainbowTableSearchService.FindRainbowTableSearchById(searchId)
	if rainbowTableSearch.ID == 0 {
		return fmt.Errorf("no rainbow table search object found for ID %d", searchId)
	}

	rainbowTable := service.rainbowTableService.FindRainbowTableById(rainbowTableSearch.RainbowTableId)
	if rainbowTable.ID == 0 {
		return fmt.Errorf("no rainbow table object found for ID %d", rainbowTableSearch.RainbowTableId)
	}

	reductionFunctionFamily := getDefaultReductionFunctionFamily(int(rainbowTable.PasswordLength), rainbowTable.CharacterSet)
	hashFunctionProvider, err := GetHashFunctionProvider(rainbowTable.HashFunction)
	if err != nil {
		return err
	}

	err = service.rainbowTableSearchService.UpdateRainbowTableSearchStatusAndSearchStarted(searchId, model.StatusStarted)
	if err != nil {
		return fmt.Errorf("failed to update status for search %d to %s: %s", searchId, model.StatusStarted, err)
	}

	indexByEndHash := make(map[string]int64)
	errorChannels := initializeErrorChannels(service.jobConfig.NumThreads)

	err = service.spawnChainGenerationThreads(
		rainbowTableSearch,
		rainbowTable,
		indexByEndHash,
		errorChannels,
		hashFunctionProvider,
	)

	searchBatches := service.getSearchBatches(indexByEndHash)
	foundChannels := initializeFoundChannels(len(searchBatches))

	for i, searchBatch := range searchBatches {
		go service.runSearchThread(
			rainbowTable.ID,
			searchBatch,
			foundChannels[i],
		)
	}

	rainbowChainGeneratorService := newChainGeneratorService(
		hashFunctionProvider.NewHashFunction(),
		reductionFunctionFamily,
	)

	plaintext := service.generatePlaintextFromFoundEndHash(
		rainbowTableSearch.Hash,
		foundChannels,
		indexByEndHash,
		rainbowChainGeneratorService,
	)

	var searchResult string

	if plaintext == "" {
		searchResult = model.StatusNotFound
	} else {
		searchResult = model.StatusFound
	}

	err = service.rainbowTableSearchService.UpdateRainbowTableSearchStatusPasswordAndSearchCompleted(
		searchId,
		searchResult,
		plaintext,
	)

	if err != nil {
		return fmt.Errorf("failed to update result of search %d to %s: %s", searchId, searchResult, err)
	}

	return nil
}
