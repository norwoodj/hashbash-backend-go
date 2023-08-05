package api

import (
	"encoding/json"
	"github.com/norwoodj/hashbash-backend-go/pkg/api_model"
	"net/http"

	"github.com/gorilla/mux"
)

func AddVersionRoutes(router *mux.Router, buildTimestamp string, gitRevision string, version string) {
	router.
		HandleFunc("/server-version.json", getVersionInfoHandler(buildTimestamp, gitRevision, version)).
		Methods("GET")
}

func getVersionInfoHandler(buildTimestamp string, gitRevision string, version string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		json.
			NewEncoder(writer).
			Encode(api_model.VersionInfo{
				BuildTimestamp: buildTimestamp,
				GitRevision:    gitRevision,
				Version:        version,
			})
	}
}
