package rainbow

import (
	"fmt"
	"github.com/norwoodj/hashbash-backend-go/pkg/dao"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
)

type TableSearchJobService struct {
	rainbowChainGeneratorService *chainGeneratorService
	rainbowChainService          dao.RainbowChainService
	rainbowTableSearchService    dao.RainbowTableSearchService
}

func (service *TableSearchJobService) generateChainIndexHashPairForHash(hash string) (int64, string) {
	return 0, ""
}

func (service *TableSearchJobService) PerformHashSearch(searchId int64) (string, error) {
	rainbowTableSearch := service.rainbowTableSearchService.FindRainbowTableSearchById(searchId)
	if rainbowTableSearch.ID == 0 {
		return "", fmt.Errorf("no rainbow table search object found for ID %d", searchId)
	}

	err := service.rainbowTableSearchService.UpdateRainbowTableSearchStatus(searchId, model.StatusStarted)
	if err != nil {
		return "", fmt.Errorf("failed to update status for search %d to %s: %s", searchId, model.StatusStarted, err)
	}

	return "", nil
}
