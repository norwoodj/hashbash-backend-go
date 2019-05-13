package api

import (
	"encoding/json"
	"fmt"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
	"github.com/norwoodj/hashbash-backend-go/pkg/mq"
	"github.com/norwoodj/hashbash-backend-go/pkg/util"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/norwoodj/hashbash-backend-go/pkg/service"
	log "github.com/sirupsen/logrus"
)

const rainbowTableDefaultChainLength = 10000
const rainbowTableDefaultCharset = "abcdefghijklmnopqrstuvwxyz"
const rainbowTableDefaultNumChains = 100000
const rainbowTableDefaultHashFunction = "MD5"
const rainbowTableDefaultPasswordLength = 8

type GenerateRainbowTableRequest struct {
	Name           string `schema:"name,required"`
	ChainLength    int64  `schema:"chainLength"`
	Charset        string `schema:"charset"`
	HashFunction   string `schema:"hashFunction"`
	NumChains      int64  `schema:"numChains"`
	PasswordLength int64  `schema:"passwordLength"`
}

func rainbowTableFromRequest(generateRequest GenerateRainbowTableRequest) model.RainbowTable {
	return model.RainbowTable{
		Name:           generateRequest.Name,
		ChainLength:    util.IntOrDefault(generateRequest.ChainLength, rainbowTableDefaultChainLength),
		CharacterSet:   util.StringOrDefault(generateRequest.Charset, rainbowTableDefaultCharset),
		HashFunction:   util.StringOrDefault(generateRequest.HashFunction, rainbowTableDefaultHashFunction),
		NumChains:      util.IntOrDefault(generateRequest.NumChains, rainbowTableDefaultNumChains),
		PasswordLength: util.IntOrDefault(generateRequest.PasswordLength, rainbowTableDefaultPasswordLength),
		Status:         model.StatusQueued,
	}
}

func AddRainbowTableRoutes(router *mux.Router, service service.RainbowTableService, producers mq.HashbashMqProducers) {
	router.
		HandleFunc("/api/rainbow-table", getListRainbowTablesHandler(service)).
		Methods("GET")

	router.
		HandleFunc("/api/rainbow-table/{rainbowTableId:[0-9]+}", getRainbowTableByIdHandler(service)).
		Methods("GET")

	router.
		HandleFunc("/api/rainbow-table/count", getCountRainbowTablesHandler(service)).
		Methods("GET")

	router.
		HandleFunc("/api/rainbow-table", getGenerateRainbowTableHandler(service, producers)).
		Methods("POST")
}

func getListRainbowTablesHandler(rainbowTableService service.RainbowTableService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		pageConfig, err := getPageConfigFromRequest(writer, request)
		if err != nil {
			return
		}

		rainbowTables := rainbowTableService.ListRainbowTables(pageConfig)
		writer.Header().Set("Content-Type", "application/json")
		json.
			NewEncoder(writer).
			Encode(rainbowTables)
	}
}

func getRainbowTableByIdHandler(rainbowTableService service.RainbowTableService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableId, err := getIdPathParamValue("rainbowTableId", writer, request, 16)
		if err != nil {
			return
		}

		rainbowTable := rainbowTableService.FindRainbowTableById(convertRainbowTableId(rainbowTableId))

		if rainbowTable.Name == "" {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		json.
			NewEncoder(writer).
			Encode(rainbowTable)
	}
}

func getCountRainbowTablesHandler(rainbowTableService service.RainbowTableService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableCount := rainbowTableService.CountRainbowTables()
		writer.Header().Set("Content-Type", "application/json")
		json.
			NewEncoder(writer).
			Encode(map[string]int64{RainbowTableCountKey: rainbowTableCount})
	}
}

func getGenerateRainbowTableHandler(
	rainbowTableService service.RainbowTableService,
	hashbashMqProducers mq.HashbashMqProducers,
) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			log.Warnf("Failed to parse generateRainbowTable form request: %s", err)
			http.Redirect(
				writer,
				request,
				fmt.Sprintf("/rainbow-tables?error=%s", url.PathEscape(err.Error())),
				http.StatusTemporaryRedirect,
			)

			return
		}

		var decoder = schema.NewDecoder()
		var generateRequest GenerateRainbowTableRequest
		err = decoder.Decode(&generateRequest, request.PostForm)

		if err != nil {
			log.Warnf("Failed to unmarshal generateRainbowTable request: %s", err)
			http.Redirect(
				writer,
				request,
				fmt.Sprintf("/rainbow-tables?error=%s", url.PathEscape(err.Error())),
				http.StatusTemporaryRedirect,
			)
		}

		rainbowTable := rainbowTableFromRequest(generateRequest)
		rainbowTableService.CreateRainbowTable(&rainbowTable)
		log.Infof("Created rainbow table %s with id %d. Publishing request for generation...", rainbowTable.Name, rainbowTable.ID)

		err = hashbashMqProducers.GenerateRainbowTableProducer.
			PublishMessage(mq.RainbowTableMessage{RainbowTableId: rainbowTable.ID})

		if err != nil {
			log.Errorf("Failed to publish generateRainbowTable request: %s", err)
			http.Redirect(
				writer,
				request,
				fmt.Sprintf("/rainbow-tables?error=%s", url.PathEscape(err.Error())),
				http.StatusTemporaryRedirect,
			)
		}

		http.Redirect(
			writer,
			request,
			"/rainbow-tables",
			http.StatusTemporaryRedirect,
		)
	}
}
