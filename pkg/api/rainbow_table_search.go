package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/norwoodj/hashbash-backend-go/pkg/service"
)

func AddRainbowTableSearchRoutes(router *mux.Router, rainbowTableSearchService service.RainbowTableSearchService) {
	router.
		HandleFunc("/api/rainbow-table/{rainbowTableId:[0-9]+}/search", getRainbowTableSearchesByIdHandler(rainbowTableSearchService)).
		Methods("GET")

	router.
		HandleFunc("/api/rainbow-table/{rainbowTableId:[0-9]+}/search/count", getCountRainbowTableSearchesHandler(rainbowTableSearchService)).
		Methods("GET")

	router.
		HandleFunc("/api/rainbow-table/{rainbowTableId:[0-9]+}/searchResults", getRainbowTableSearchResultsHandler(rainbowTableSearchService)).
		Methods("GET")

	router.
		HandleFunc("/api/rainbow-table/search/{searchId:[0-9]+}", getRainbowTableSearchByIdHandler(rainbowTableSearchService)).
		Methods("GET")
}

func getRainbowTableSearchesByIdHandler(rainbowTableSearchService service.RainbowTableSearchService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableId, err := getIdPathParamValue("rainbowTableId", writer, request, 16)
		if err != nil {
			return
		}

		pageConfig, err := getPageConfigFromRequest(writer, request)
		if err != nil {
			return
		}

		includeNotFound := getIncludeNotFoundQueryParam(request.URL.Query())
		rainbowTableSearches := rainbowTableSearchService.ListSearchesByRainbowTableId(convertRainbowTableId(rainbowTableId), includeNotFound, pageConfig)
		writer.Header().Set("Content-Type", "application/json")
		json.
			NewEncoder(writer).
			Encode(rainbowTableSearches)
	}
}

func getCountRainbowTableSearchesHandler(rainbowTableSearchService service.RainbowTableSearchService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableId, err := getIdPathParamValue("rainbowTableId", writer, request, 16)
		if err != nil {
			return
		}

		includeNotFound := getIncludeNotFoundQueryParam(request.URL.Query())
		rainbowTableSearchCount := rainbowTableSearchService.CountRainbowTableSearches(convertRainbowTableId(rainbowTableId), includeNotFound)
		writer.Header().Set("Content-Type", "application/json")
		json.
			NewEncoder(writer).
			Encode(map[string]int64{RainbowTableSearchCountKey: rainbowTableSearchCount})
	}
}

func getRainbowTableSearchResultsHandler(rainbowTableSearchService service.RainbowTableSearchService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableId, err := getIdPathParamValue("rainbowTableId", writer, request, 16)
		if err != nil {
			return
		}

		searchResultResponse := rainbowTableSearchService.GetRainbowTableSearchResults(convertRainbowTableId(rainbowTableId))
		writer.Header().Set("Content-Type", "application/json")
		json.
			NewEncoder(writer).
			Encode(searchResultResponse)
	}
}

func getRainbowTableSearchByIdHandler(rainbowTableSearchService service.RainbowTableSearchService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableSearchId, err := getIdPathParamValue("searchId", writer, request, 64)
		if err != nil {
			return
		}

		rainbowTableSearches := rainbowTableSearchService.FindRainbowTableSearchById(rainbowTableSearchId.(int64))
		writer.Header().Set("Content-Type", "application/json")
		json.
			NewEncoder(writer).
			Encode(rainbowTableSearches)
	}
}
