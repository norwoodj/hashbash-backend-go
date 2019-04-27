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
	rainbowTableRouter := router.
		Path("/api/rainbow-table").
		Subrouter()

	rainbowTableRouter.
		HandleFunc("", getListRainbowTablesHandler(db)).
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

		json.NewEncoder(writer).
			Encode(rainbowTables)
	}
}
