package service

import (
	"github.com/jinzhu/gorm"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
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
	ListRainbowTableSearches(rainbowTableId int16, includeNotFound bool, pageConfig PageConfig) []model.RainbowTableSearch
	GetRainbowTableSearchResults(rainbowTableId int16) RainbowTableSearchResultSummary
}

type MySQLRainbowTableSearchService struct {
	DatabaseClient *gorm.DB
}

func NewRainbowTableSearchService(db *gorm.DB) RainbowTableSearchService {
	return MySQLRainbowTableSearchService{DatabaseClient: db}
}

func (service MySQLRainbowTableSearchService) CountRainbowTableSearches(rainbowTableId int16, includeNotFound bool) int64 {
	var rainbowTableSearchCount int64
	query := service.DatabaseClient.
		Model(&model.RainbowTableSearch{}).
		Where("rainbowTableId = ?", rainbowTableId)

	if !includeNotFound {
		query = query.Where("status != ?", model.SearchNotFound)
	}

	query.Count(&rainbowTableSearchCount)
	return rainbowTableSearchCount
}

func (service MySQLRainbowTableSearchService) ListRainbowTableSearches(
	rainbowTableId int16,
	includeNotFound bool,
	pageConfig PageConfig,
) []model.RainbowTableSearch {
	rainbowTableSearches := make([]model.RainbowTableSearch, 0)
	query := applyPaging(service.DatabaseClient, pageConfig).
		Where("rainbowTableId = ?", rainbowTableId)

	if !includeNotFound {
		query = query.Where("status != ?", model.SearchNotFound)
	}

	query.Find(&rainbowTableSearches)
	return rainbowTableSearches
}

func (service MySQLRainbowTableSearchService) GetRainbowTableSearchResults(rainbowTableId int16) RainbowTableSearchResultSummary {
	searchResults := make([]RainbowTableSearchResults, 0)
	service.DatabaseClient.
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
