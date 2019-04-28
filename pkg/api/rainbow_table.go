package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/norwoodj/hashbash-backend-go/pkg/service"
)

func AddRainbowTableRoutes(router *mux.Router, service service.RainbowTableService) {
	router.
		HandleFunc("/api/rainbow-table", getListRainbowTablesHandler(service)).
		Methods("GET")

	router.
		HandleFunc("/api/rainbow-table/{id:[0-9]+}", getRainbowTableByIdHandler(service)).
		Methods("GET")

	router.
		HandleFunc("/api/rainbow-table/count", getCountRainbowTablesHandler(service)).
		Methods("GET")
}

func getListRainbowTablesHandler(service service.RainbowTableService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		pageConfig, err := getPageConfigFromRequest(writer, request)
		if err != nil {
			return
		}

		rainbowTables := service.ListRainbowTables(pageConfig)
		json.
			NewEncoder(writer).
			Encode(rainbowTables)
	}
}

func getRainbowTableByIdHandler(service service.RainbowTableService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableId, err := getIdPathParamValue(writer, request)
		if err != nil {
			return
		}

		rainbowTable := service.FindRainbowTableById(rainbowTableId)

		if rainbowTable.Name == "" {
			writer.WriteHeader(404)
			return
		}

		json.
			NewEncoder(writer).
			Encode(rainbowTable)
	}
}

func getCountRainbowTablesHandler(service service.RainbowTableService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableCount := service.CountRainbowTables()
		json.
			NewEncoder(writer).
			Encode(map[string]int64{RainbowTableCountKey: rainbowTableCount})
	}
}
