package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
	"time"
)

type RainbowTableSearchResults struct {
	Count        int64  `json:"count"`
	SearchStatus string `json:"searchStatus"`
}

type RainbowTableSearchResultSummary struct {
	FoundSearches int64 `json:"foundSearches"`
	TotalSearches int64 `json:"totalSearches"`
}

type RainbowTableSearchService interface {
	CountRainbowTableSearches(rainbowTableId int16, includeNotFound bool) int64
	CreateRainbowTableSearch(rainbowTableId int16, hash string) (model.RainbowTableSearch, error)
	ListSearchesByRainbowTableId(rainbowTableId int16, includeNotFound bool, pageConfig PageConfig) []model.RainbowTableSearch
	GetRainbowTableSearchResults(rainbowTableId int16) RainbowTableSearchResultSummary
	FindRainbowTableSearchById(int64) model.RainbowTableSearch
	UpdateRainbowTableSearchStatus(int64, string) error
	UpdateRainbowTableSearchStatusAndSearchStarted(int64, string) error
	UpdateRainbowTableSearchStatusPasswordAndSearchCompleted(int64, string, string) error
}

type MySQLRainbowTableSearchService struct {
	databaseClient *gorm.DB
}

func NewRainbowTableSearchService(db *gorm.DB) RainbowTableSearchService {
	return &MySQLRainbowTableSearchService{databaseClient: db}
}

func (service *MySQLRainbowTableSearchService) CountRainbowTableSearches(rainbowTableId int16, includeNotFound bool) int64 {
	var rainbowTableSearchCount int64
	query := service.databaseClient.
		Model(&model.RainbowTableSearch{}).
		Where("rainbowTableId = ?", rainbowTableId)

	if !includeNotFound {
		query = query.Where("status != ?", model.SearchNotFound)
	}

	query.Count(&rainbowTableSearchCount)
	return rainbowTableSearchCount
}

func (service *MySQLRainbowTableSearchService) CreateRainbowTableSearch(rainbowTableId int16, hash string) (model.RainbowTableSearch, error) {
	var rainbowTable model.RainbowTable
	service.databaseClient.
		Where("id = ?", rainbowTableId).
		First(&rainbowTable)

	if rainbowTable.Name == "" {
		return model.RainbowTableSearch{}, RainbowTableNotExistsError{ID: rainbowTableId}
	}

	if !isValidHash(rainbowTable.HashFunction, hash) {
		return model.RainbowTableSearch{}, InvalidHashError{Hash: hash, HashFunctionName: rainbowTable.HashFunction}
	}

	rainbowTableSearch := model.RainbowTableSearch{
		RainbowTableId: rainbowTableId,
		Hash:           hash,
		Status:         model.StatusQueued,
	}

	err := service.databaseClient.
		Save(&rainbowTableSearch).
		Error

	return rainbowTableSearch, err
}

func (service *MySQLRainbowTableSearchService) ListSearchesByRainbowTableId(
	rainbowTableId int16,
	includeNotFound bool,
	pageConfig PageConfig,
) []model.RainbowTableSearch {
	rainbowTableSearches := make([]model.RainbowTableSearch, 0)
	query := applyPaging(service.databaseClient, pageConfig).
		Where("rainbowTableId = ?", rainbowTableId)

	if !includeNotFound {
		query = query.Where("status != ?", model.SearchNotFound)
	}

	query.Find(&rainbowTableSearches)
	return rainbowTableSearches
}

func (service *MySQLRainbowTableSearchService) GetRainbowTableSearchResults(rainbowTableId int16) RainbowTableSearchResultSummary {
	searchResults := make([]RainbowTableSearchResults, 0)
	service.databaseClient.
		Model(&model.RainbowTableSearch{}).
		Select("status AS searchStatus, COUNT(*) AS count").
		Where("rainbowTableId = ? and status IN (?)", rainbowTableId, []string{model.SearchFound, model.SearchNotFound}).
		Group("searchStatus").
		Scan(&searchResults)

	var searchResultSummary RainbowTableSearchResultSummary

	for _, result := range searchResults {
		if result.SearchStatus == model.SearchFound {
			searchResultSummary.FoundSearches += result.Count
		}

		searchResultSummary.TotalSearches += result.Count
	}

	return searchResultSummary
}

func (service *MySQLRainbowTableSearchService) FindRainbowTableSearchById(searchId int64) model.RainbowTableSearch {
	var rainbowTableSearch model.RainbowTableSearch

	service.databaseClient.
		Where("id = ?", searchId).
		First(&rainbowTableSearch)

	return rainbowTableSearch
}

func (service *MySQLRainbowTableSearchService) UpdateRainbowTableSearchStatus(searchId int64, status string) error {
	return service.databaseClient.
		Model(&model.RainbowTableSearch{ID: searchId}).
		Update("status", status).
		Error
}

func (service *MySQLRainbowTableSearchService) UpdateRainbowTableSearchStatusAndSearchStarted(searchId int64, status string) error {
	return service.databaseClient.
		Model(&model.RainbowTableSearch{ID: searchId}).
		Updates(map[string]interface{}{
			"status":        status,
			"searchStarted": time.Now(),
		}).
		Error
}

func (service *MySQLRainbowTableSearchService) UpdateRainbowTableSearchStatusPasswordAndSearchCompleted(
	searchId int64,
	status string,
	password string,
) error {
	return service.databaseClient.
		Model(&model.RainbowTableSearch{ID: searchId}).
		Updates(map[string]interface{}{
			"status":          status,
			"password":        password,
			"searchCompleted": time.Now(),
		}).
		Error
}
