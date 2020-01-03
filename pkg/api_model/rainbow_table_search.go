package api_model

import (
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
	"time"
)

type RainbowTableSearchRequest struct {
	Hash string `json:"hash"`
}

type RainbowTableSearchResponse struct {
	Hash     string `json:"hash"`
	SearchId int64  `json:"searchId"`
	Status   string `json:"status"`
}

type RainbowTableSearch struct {
	ID              int64      `json:"id"`
	RainbowTableId  int16      `json:"rainbowTableId"`
	Hash            string     `json:"hash"`
	Status          string     `json:"status"`
	Password        string     `json:"password"`
	SearchStarted   *time.Time `json:"searchStarted"`
	SearchCompleted *time.Time `json:"searchCompleted"`
	SearchTime      *float64   `json:"searchTime"`
	CreatedAt       *time.Time `json:"created"`
}

func ConvertRainbowTableSearchToApiModel(rainbowTableSearch model.RainbowTableSearch) RainbowTableSearch {
	var searchTimeResult *float64
	if rainbowTableSearch.SearchCompleted != nil {
		searchTime := rainbowTableSearch.SearchCompleted.Sub(*rainbowTableSearch.SearchStarted).Seconds()
		searchTimeResult = &searchTime
	}

	return RainbowTableSearch{
		ID:              rainbowTableSearch.ID,
		RainbowTableId:  rainbowTableSearch.RainbowTableId,
		Hash:            rainbowTableSearch.Hash,
		Status:          rainbowTableSearch.Status,
		Password:        rainbowTableSearch.Password,
		SearchStarted:   rainbowTableSearch.SearchStarted,
		SearchCompleted: rainbowTableSearch.SearchCompleted,
		SearchTime:      searchTimeResult,
		CreatedAt:       rainbowTableSearch.CreatedAt,
	}
}

func ConvertRainbowTableSearchesToApiModels(rainbowTableSearches []model.RainbowTableSearch) []RainbowTableSearch {
	result := make([]RainbowTableSearch, 0)
	for _, s := range rainbowTableSearches {
		result = append(result, ConvertRainbowTableSearchToApiModel(s))
	}

	return result
}
