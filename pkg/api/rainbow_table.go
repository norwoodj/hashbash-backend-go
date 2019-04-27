package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/norwoodj/hashbash-backend-go/pkg/database"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
)

func AddRainbowTableRoutes(router *mux.Router, db *gorm.DB) {
	// GET /api/rainbow-table
	router.
		HandleFunc("/api/rainbow-table", getListRainbowTablesHandler(db)).
		Methods("GET")

	// GET /api/rainbow-table/{id}
	router.
		HandleFunc("/api/rainbow-table/{id:[0-9]+}", getRainbowTableByIdHandler(db)).
		Methods("GET")

	// GET /api/rainbow-table/count
	router.
		HandleFunc("/api/rainbow-table/count", getCountRainbowTablesHandler(db)).
		Methods("GET")
}

func getListRainbowTablesHandler(db *gorm.DB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		pageConfig, err := getPageConfigFromRequest(writer, request)
		if err != nil {
			return
		}

		rainbowTables := make([]model.RainbowTable, 0)
		database.ApplyPaging(db, pageConfig).
			Find(&rainbowTables)

		json.
			NewEncoder(writer).
			Encode(rainbowTables)
	}
}

func getRainbowTableByIdHandler(db *gorm.DB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableId, err := getIdPathParamValue(writer, request)
		if err != nil {
			return
		}

		var rainbowTable model.RainbowTable
		db.
			Where("id = ?", rainbowTableId).
			First(&rainbowTable)

		if rainbowTable.Name == "" {
			writer.WriteHeader(404)
			return
		}

		json.
			NewEncoder(writer).
			Encode(rainbowTable)
	}
}

func getCountRainbowTablesHandler(db *gorm.DB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var rainbowTableCount int64
		db.
			Model(&model.RainbowTable{}).
			Count(&rainbowTableCount)

		json.
			NewEncoder(writer).
			Encode(map[string]int64{"rainbowTableCount": rainbowTableCount})
	}
}
