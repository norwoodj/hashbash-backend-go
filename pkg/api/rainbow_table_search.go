package api

import (
	"encoding/json"
	"github.com/norwoodj/hashbash-backend-go/pkg/service"
	"net/http"

	"github.com/gorilla/mux"
)

func AddRainbowTableSearchRoutes(router *mux.Router, rainbowTableSearchService service.RainbowTableSearchService) {
	router.
		HandleFunc("/api/rainbow-table/{id:[0-9]+}/search", getRainbowTableSearchesByIdHandler(rainbowTableSearchService)).
		Methods("GET")

	router.
		HandleFunc("/api/rainbow-table/{id:[0-9]+}/search/count", getCountRainbowTableSearchesHandler(rainbowTableSearchService)).
		Methods("GET")

	router.
		HandleFunc("/api/rainbow-table/{id:[0-9]+}/searchResults", getRainbowTableSearchResultsHandler(rainbowTableSearchService)).
		Methods("GET")
}

func getRainbowTableSearchesByIdHandler(service service.RainbowTableSearchService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableId, err := getIdPathParamValue(writer, request)
		if err != nil {
			return
		}

		pageConfig, err := getPageConfigFromRequest(writer, request)
		if err != nil {
			return
		}

		includeNotFound := getIncludeNotFoundQueryParam(request.URL.Query())
		rainbowTableSearches := service.ListRainbowTableSearches(rainbowTableId, includeNotFound, pageConfig)
		writer.Header().Set("Content-Type", "application/json")
		json.
			NewEncoder(writer).
			Encode(rainbowTableSearches)
	}
}

func getCountRainbowTableSearchesHandler(service service.RainbowTableSearchService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableId, err := getIdPathParamValue(writer, request)
		if err != nil {
			return
		}

		includeNotFound := getIncludeNotFoundQueryParam(request.URL.Query())
		rainbowTableSearchCount := service.CountRainbowTableSearches(rainbowTableId, includeNotFound)
		writer.Header().Set("Content-Type", "application/json")
		json.
			NewEncoder(writer).
			Encode(map[string]int64{RainbowTableSearchCountKey: rainbowTableSearchCount})
	}
}

func getRainbowTableSearchResultsHandler(service service.RainbowTableSearchService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableId, err := getIdPathParamValue(writer, request)
		if err != nil {
			return
		}

		searchResultResponse := service.GetRainbowTableSearchResults(rainbowTableId)
		writer.Header().Set("Content-Type", "application/json")
		json.
			NewEncoder(writer).
			Encode(searchResultResponse)
	}
}
