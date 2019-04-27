package api

import (
	"encoding/json"
	"github.com/norwoodj/hashbash-backend-go/pkg/database"
	"github.com/norwoodj/hashbash-backend-go/pkg/util"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
)

type RainbowTableSearchResults struct {
	Count        int64  `json:"count"`
	SearchStatus string `json:"searchStatus"`
}

type RainbowTableSearchResultResponse struct {
	FoundSearches int64 `json:"foundSearches"`
	TotalSearches int64 `json:"totalSearches"`
}

func AddRainbowTableSearchRoutes(router *mux.Router, db *gorm.DB) {
	// GET /api/rainbow-table/{id}/search
	router.
		HandleFunc("/api/rainbow-table/{id:[0-9]+}/search", getRainbowTableSearchesByIdHandler(db)).
		Methods("GET")

	// GET /api/rainbow-table/{id}/search/count
	router.
		HandleFunc("/api/rainbow-table/{id:[0-9]+}/search/count", getCountRainbowTableSearchesHandler(db)).
		Methods("GET")

	// GET /api/rainbow-table/{rainbowTableId}/searchResults
	router.
		HandleFunc("/api/rainbow-table/{id:[0-9]+}/searchResults", getRainbowTableSearchResultsHandler(db)).
		Methods("GET")
}

func getRainbowTableSearchesByIdHandler(db *gorm.DB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableId, err := getIdPathParamValue(writer, request)
		if err != nil {
			return
		}

		pageConfig, err := getPageConfigFromRequest(writer, request)
		if err != nil {
			return
		}

		rainbowTableSearches := make([]model.RainbowTableSearch, 0)
		database.ApplyPaging(db, pageConfig).
			Where("rainbowTableId = ?", rainbowTableId).
			Find(&rainbowTableSearches)

		json.
			NewEncoder(writer).
			Encode(rainbowTableSearches)
	}
}

func getCountRainbowTableSearchesHandler(db *gorm.DB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableId, err := getIdPathParamValue(writer, request)
		if err != nil {
			return
		}

		var rainbowTableSearchCount int64
		db.
			Model(&model.RainbowTableSearch{}).
			Where("rainbowTableId = ?", rainbowTableId).
			Count(&rainbowTableSearchCount)

		json.
			NewEncoder(writer).
			Encode(map[string]int64{"searchCount": rainbowTableSearchCount})
	}
}

func getRainbowTableSearchResultsHandler(db *gorm.DB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableId, err := getIdPathParamValue(writer, request)
		if err != nil {
			return
		}

		searchResults := make([]RainbowTableSearchResults, 0)
		db.
			Model(&model.RainbowTableSearch{}).
			Select("status AS searchStatus, COUNT(*) AS count").
			Where("rainbowTableId = ? and status IN (?)", rainbowTableId, []string{util.Found, util.NotFound}).
			Group("searchStatus").
			Scan(&searchResults)

		var searchResultResponse RainbowTableSearchResultResponse

		for _, result := range searchResults {
			if result.SearchStatus == util.Found {
				searchResultResponse.FoundSearches += result.Count
			}

			searchResultResponse.TotalSearches += result.Count
		}

		json.
			NewEncoder(writer).
			Encode(searchResultResponse)
	}
}
