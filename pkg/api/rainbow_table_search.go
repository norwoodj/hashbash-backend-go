package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/norwoodj/hashbash-backend-go/pkg/api_model"
	"github.com/norwoodj/hashbash-backend-go/pkg/dao"
	"github.com/norwoodj/hashbash-backend-go/pkg/rabbitmq"
	log "github.com/sirupsen/logrus"
)

func AddRainbowTableSearchRoutes(
	router *mux.Router,
	rainbowTableSearchService dao.RainbowTableSearchService,
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

func getRainbowTableSearchesByIdHandler(rainbowTableSearchService dao.RainbowTableSearchService) func(writer http.ResponseWriter, request *http.Request) {
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
			Encode(api_model.ConvertRainbowTableSearchesToApiModels(rainbowTableSearches))
	}
}

func getCountRainbowTableSearchesHandler(rainbowTableSearchService dao.RainbowTableSearchService) func(writer http.ResponseWriter, request *http.Request) {
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

func getRainbowTableSearchResultsHandler(rainbowTableSearchService dao.RainbowTableSearchService) func(writer http.ResponseWriter, request *http.Request) {
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

func getRainbowTableSearchByIdHandler(rainbowTableSearchService dao.RainbowTableSearchService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableSearchId, err := getIdPathParamValue("searchId", writer, request, 64)
		if err != nil {
			return
		}

		rainbowTableSearch := rainbowTableSearchService.FindRainbowTableSearchById(rainbowTableSearchId)
		writer.Header().Set("Content-Type", "application/json")
		json.
			NewEncoder(writer).
			Encode(api_model.ConvertRainbowTableSearchToApiModel(rainbowTableSearch))
	}
}

func createRainbowTableSearchByIdHandler(
	rainbowTableSearchService dao.RainbowTableSearchService,
	hashbashMqProducers rabbitmq.HashbashMqProducers,
) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableId, err := getIdPathParamValue("rainbowTableId", writer, request, 16)
		if err != nil {
			return
		}

		hash := getStringQueryParamValue(
			request.URL.Query(),
			"hash",
			writer,
		)

		rainbowTableSearch, err := rainbowTableSearchService.CreateRainbowTableSearch(
			int16(rainbowTableId),
			hash,
		)

		if err != nil {
			if dao.IsRainbowTableNotExistsError(err) || dao.IsInvalidHashError(err) {
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
			Hash:           hash,
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

		response := api_model.RainbowTableSearchResponse{
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
