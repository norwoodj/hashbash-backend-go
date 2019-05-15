package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/norwoodj/hashbash-backend-go/pkg/rabbitmq"
	"github.com/norwoodj/hashbash-backend-go/pkg/service"
	log "github.com/sirupsen/logrus"
)

func AddRainbowTableSearchRoutes(
	router *mux.Router,
	rainbowTableSearchService service.RainbowTableSearchService,
	hashbashMqProducers rabbitmq.HashbashMqProducers,
) {
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

	router.
		HandleFunc("/api/rainbow-table/{rainbowTableId:[0-9]+}/search", createRainbowTableSearchByIdHandler(rainbowTableSearchService, hashbashMqProducers)).
		Methods("POST")
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

		rainbowTableSearches := rainbowTableSearchService.FindRainbowTableSearchById(rainbowTableSearchId)
		writer.Header().Set("Content-Type", "application/json")
		json.
			NewEncoder(writer).
			Encode(rainbowTableSearches)
	}
}

func createRainbowTableSearchByIdHandler(
	rainbowTableSearchService service.RainbowTableSearchService,
	hashbashMqProducers rabbitmq.HashbashMqProducers,
) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableId, err := getIdPathParamValue("rainbowTableId", writer, request, 16)
		if err != nil {
			return
		}

		requestBody, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Warnf("Failed to read request body: %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		var searchRequest rainbowTableSearchRequest
		err = json.Unmarshal(requestBody, &searchRequest)
		if err != nil {
			log.Warnf("Failed to unmarshal rainbow table search request: %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		rainbowTableSearch, err := rainbowTableSearchService.CreateRainbowTableSearch(
			int16(rainbowTableId),
			searchRequest.Hash,
		)

		if err != nil {
			if service.IsRainbowTableNotExistsError(err) || service.IsInvalidHashError(err) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusBadRequest)
				json.
					NewEncoder(writer).
					Encode(map[string]string{"error": err.Error()})

			} else {
				writer.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		searchRequestMessage := rabbitmq.RainbowTableSearchRequestMessage{
			Hash:           searchRequest.Hash,
			RainbowTableId: int16(rainbowTableId),
			SearchId:       rainbowTableSearch.ID,
		}

		err = hashbashMqProducers.SearchRainbowTableProducer.
			PublishMessage(searchRequestMessage)

		if err != nil {
			log.Errorf("Unknown error occurred publishing search request for rainbow table %d: %s", rainbowTableId, err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := rainbowTableSearchResponse{
			Hash:     rainbowTableSearch.Hash,
			SearchId: rainbowTableSearch.ID,
			Status:   rainbowTableSearch.Status,
		}

		writer.Header().Set("Content-Type", "application/json")
		json.
			NewEncoder(writer).
			Encode(response)
	}
}
