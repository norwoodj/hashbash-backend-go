package dao

import (
	"encoding/hex"
	"github.com/jinzhu/gorm"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
	"time"
)

type RainbowTableSearchResults struct {
	Count        int64  `gorm:"column:count"json:"count"`
	SearchStatus string `gorm:"column:status"json:"searchStatus"`
}

type RainbowTableSearchResultSummary struct {
	FoundSearches int64 `json:"foundSearches"`
	TotalSearches int64 `json:"totalSearches"`
}

type RainbowTableSearchService interface {
	CountRainbowTableSearches(rainbowTableId int16, includeNotFound bool) (int64, error)
	CreateRainbowTableSearch(rainbowTableId int16, hash string) (model.RainbowTableSearch, error)
	ListSearchesByRainbowTableId(rainbowTableId int16, includeNotFound bool, pageConfig PageConfig) ([]model.RainbowTableSearch, error)
	GetRainbowTableSearchResults(rainbowTableId int16) (RainbowTableSearchResultSummary, error)
	FindRainbowTableSearchById(int64) (model.RainbowTableSearch, error)
	UpdateRainbowTableSearchStatus(int64, string) error
	UpdateRainbowTableSearchStatusAndSearchStarted(int64, string) error
	UpdateRainbowTableSearchStatusPasswordAndSearchCompleted(int64, string, string) error
}

type DbRainbowTableSearchService struct {
	databaseClient *gorm.DB
}

func NewRainbowTableSearchService(db *gorm.DB) RainbowTableSearchService {
	return &DbRainbowTableSearchService{databaseClient: db}
}

func (service *DbRainbowTableSearchService) CountRainbowTableSearches(rainbowTableId int16, includeNotFound bool) (int64, error) {
	var rainbowTableSearchCount int64
	query := service.databaseClient.
		Model(&model.RainbowTableSearch{}).
		Where("rainbow_table_id = ?", rainbowTableId)

	if !includeNotFound {
		query = query.Where("status != ?", model.SearchNotFound)
	}

	err := query.Count(&rainbowTableSearchCount).Error
	return rainbowTableSearchCount, err
}

func (service *DbRainbowTableSearchService) CreateRainbowTableSearch(rainbowTableId int16, hash string) (model.RainbowTableSearch, error) {
	var rainbowTable model.RainbowTable
	service.databaseClient.
		Where("id = ?", rainbowTableId).
		First(&rainbowTable)

	if rainbowTable.Name == "" {
		return model.RainbowTableSearch{}, RainbowTableNotExistsError{ID: rainbowTableId}
	}

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return model.RainbowTableSearch{}, InvalidHashError{Hash: hash, HashFunctionName: rainbowTable.HashFunction}
	}

	rainbowTableSearch := model.RainbowTableSearch{
		RainbowTableId: rainbowTableId,
		Hash:           hashBytes,
		Status:         model.StatusQueued,
	}

	err = service.databaseClient.
		Save(&rainbowTableSearch).
		Error

	return rainbowTableSearch, err
}

func (service *DbRainbowTableSearchService) ListSearchesByRainbowTableId(
	rainbowTableId int16,
	includeNotFound bool,
	pageConfig PageConfig,
) ([]model.RainbowTableSearch, error) {
	rainbowTableSearches := make([]model.RainbowTableSearch, 0)
	query := applyPaging(service.databaseClient, pageConfig).
		Where("rainbow_table_id = ?", rainbowTableId)

	if !includeNotFound {
		query = query.Where("status != ?", model.SearchNotFound)
	}

	err := query.Find(&rainbowTableSearches).Error
	return rainbowTableSearches, err
}

func (service *DbRainbowTableSearchService) GetRainbowTableSearchResults(rainbowTableId int16) (RainbowTableSearchResultSummary, error) {
	searchResults := make([]RainbowTableSearchResults, 0)
	err := service.databaseClient.
		Model(&model.RainbowTableSearch{}).
		Select("status, COUNT(*) AS count").
		Where("rainbow_table_id = ? and status IN (?)", rainbowTableId, []string{model.SearchFound, model.SearchNotFound}).
		Group("status").
		Scan(&searchResults).
		Error

	if err != nil {
		return RainbowTableSearchResultSummary{}, err
	}

	var searchResultSummary RainbowTableSearchResultSummary

	for _, result := range searchResults {
		if result.SearchStatus == model.SearchFound {
			searchResultSummary.FoundSearches += result.Count
		}

		searchResultSummary.TotalSearches += result.Count
	}

	return searchResultSummary, nil
}

func (service *DbRainbowTableSearchService) FindRainbowTableSearchById(searchId int64) (model.RainbowTableSearch, error) {
	var rainbowTableSearch model.RainbowTableSearch

	err := service.databaseClient.
		Where("id = ?", searchId).
		First(&rainbowTableSearch).
		Error

	return rainbowTableSearch, err
}

func (service *DbRainbowTableSearchService) UpdateRainbowTableSearchStatus(searchId int64, status string) error {
	return service.databaseClient.
		Model(&model.RainbowTableSearch{ID: searchId}).
		Update("status", status).
		Error
}

func (service *DbRainbowTableSearchService) UpdateRainbowTableSearchStatusAndSearchStarted(searchId int64, status string) error {
	return service.databaseClient.
		Model(&model.RainbowTableSearch{ID: searchId}).
		Updates(map[string]interface{}{
			"status":         status,
			"search_started": time.Now(),
		}).
		Error
}

func (service *DbRainbowTableSearchService) UpdateRainbowTableSearchStatusPasswordAndSearchCompleted(
	searchId int64,
	status string,
	password string,
) error {
	return service.databaseClient.
		Model(&model.RainbowTableSearch{ID: searchId}).
		Updates(map[string]interface{}{
			"status":           status,
			"password":         password,
			"search_completed": time.Now(),
		}).
		Error
}
