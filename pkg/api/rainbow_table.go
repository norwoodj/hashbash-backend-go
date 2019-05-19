package api

import (
	"encoding/json"
	"fmt"
	"github.com/norwoodj/hashbash-backend-go/pkg/api_model"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
	"github.com/norwoodj/hashbash-backend-go/pkg/rabbitmq"
	"github.com/norwoodj/hashbash-backend-go/pkg/util"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/norwoodj/hashbash-backend-go/pkg/dao"
	log "github.com/sirupsen/logrus"
)

const rainbowTableDefaultChainLength = 10000
const rainbowTableDefaultCharset = "abcdefghijklmnopqrstuvwxyz"
const rainbowTableDefaultNumChains = 100000
const rainbowTableDefaultHashFunction = "MD5"
const rainbowTableDefaultPasswordLength = 8

type GenerateRainbowTableRequest struct {
	Name           string `json:"name,required"schema:"name,required"`
	ChainLength    int64  `json:"chainLength"schema:"chainLength"`
	Charset        string `json:"charset"schema:"charset"`
	HashFunction   string `json:"hashFunction"schema:"hashFunction"`
	NumChains      int64  `json:"numChains"schema:"numChains"`
	PasswordLength int64  `json:"passwordLength"schema:"passwordLength"`
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

func AddRainbowTableRoutes(router *mux.Router, service dao.RainbowTableService, producers rabbitmq.HashbashMqProducers) {
	router.
		HandleFunc("/api/rainbow-table", getListRainbowTablesHandler(service)).
		Methods("GET")

	router.
		HandleFunc("/api/rainbow-table", getGenerateRainbowTableFormHandler(service, producers)).
		Headers("Content-Type", "application/x-www-form-urlencoded").
		Methods("POST")

	router.
		HandleFunc("/api/rainbow-table", getGenerateRainbowTableJsonHandler(service, producers)).
		Headers("Content-Type", "application/json").
		Methods("POST")

	router.
		HandleFunc("/api/rainbow-table/{rainbowTableId:[0-9]+}", getRainbowTableByIdHandler(service)).
		Methods("GET")

	router.
		HandleFunc("/api/rainbow-table/{rainbowTableId:[0-9]+}", deleteRainbowTableByIdHandler(service, producers)).
		Methods("DELETE")

	router.
		HandleFunc("/api/rainbow-table/count", getCountRainbowTablesHandler(service)).
		Methods("GET")
}

func getListRainbowTablesHandler(rainbowTableService dao.RainbowTableService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		pageConfig, err := getPageConfigFromRequest(writer, request)
		if err != nil {
			return
		}

		rainbowTables := rainbowTableService.ListRainbowTables(pageConfig)
		writer.Header().Set("Content-Type", "application/json")
		json.
			NewEncoder(writer).
			Encode(api_model.ConvertRainbowTablesToApiModels(rainbowTables))
	}
}

func getRainbowTableByIdHandler(rainbowTableService dao.RainbowTableService) func(writer http.ResponseWriter, request *http.Request) {
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
			Encode(api_model.ConvertRainbowTableToApiModel(rainbowTable))
	}
}

func getCountRainbowTablesHandler(rainbowTableService dao.RainbowTableService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rainbowTableCount := rainbowTableService.CountRainbowTables()
		writer.Header().Set("Content-Type", "application/json")
		json.
			NewEncoder(writer).
			Encode(map[string]int64{RainbowTableCountKey: rainbowTableCount})
	}
}

func handleGenerateRainbowTable(
	rainbowTableService dao.RainbowTableService,
	hashbashMqProducers rabbitmq.HashbashMqProducers,
	generateRequest GenerateRainbowTableRequest,
) (model.RainbowTable, error) {
	rainbowTable := rainbowTableFromRequest(generateRequest)
	_, err := rainbowTableService.CreateRainbowTable(&rainbowTable)

	if err != nil {
		return rainbowTable, err
	}

	log.Infof("Created rainbow table %s with id %d. Publishing request for generation...", rainbowTable.Name, rainbowTable.ID)
	err = hashbashMqProducers.GenerateRainbowTableProducer.
		PublishMessage(rabbitmq.RainbowTableIdMessage{RainbowTableId: rainbowTable.ID})

	return rainbowTable, err
}

func getGenerateRainbowTableFormHandler(
	rainbowTableService dao.RainbowTableService,
	hashbashMqProducers rabbitmq.HashbashMqProducers,
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

			return
		}

		_, err = handleGenerateRainbowTable(
			rainbowTableService,
			hashbashMqProducers,
			generateRequest,
		)

		if err != nil {
			log.Errorf("Failed to publish generateRainbowTable request: %s", err)
			http.Redirect(
				writer,
				request,
				fmt.Sprintf("/rainbow-tables?error=%s", url.PathEscape(err.Error())),
				http.StatusTemporaryRedirect,
			)

			return
		}

		http.Redirect(
			writer,
			request,
			"/rainbow-tables",
			http.StatusTemporaryRedirect,
		)
	}
}

func getGenerateRainbowTableJsonHandler(
	rainbowTableService dao.RainbowTableService,
	hashbashMqProducers rabbitmq.HashbashMqProducers,
) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestBody, err := ioutil.ReadAll(request.Body)

		if err != nil {
			log.Warnf("Failed to read request body: %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		var generateRequest GenerateRainbowTableRequest
		err = json.Unmarshal(requestBody, &generateRequest)

		if err != nil {
			log.Warnf("Failed to unmarshal generateRainbowTable request: %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if generateRequest.Name == "" {
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).
				Encode(map[string]string{"error": "'name' is required for generate rainbow table requests"})
			return
		}

		rainbowTable, err := handleGenerateRainbowTable(
			rainbowTableService,
			hashbashMqProducers,
			generateRequest,
		)

		if err != nil {
			if dao.IsRainbowTableExistsError(err) {
				writer.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(writer).
					Encode(map[string]string{"error": err.Error()})
				return
			}

			log.Errorf("Failed to publish generateRainbowTable request: %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Location", fmt.Sprintf("/api/rainbow-table/%d", rainbowTable.ID))
		writer.WriteHeader(http.StatusCreated)
	}
}

func deleteRainbowTableByIdHandler(
	rainbowTableService dao.RainbowTableService,
	hashbashMqProducers rabbitmq.HashbashMqProducers,
) func(writer http.ResponseWriter, request *http.Request) {
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

		err = hashbashMqProducers.DeleteRainbowTableProducer.
			PublishMessage(rabbitmq.RainbowTableIdMessage{RainbowTableId: rainbowTable.ID})

		if err != nil {
			log.Errorf("Failed to publish deleteRainbowTable request: %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Location", fmt.Sprintf("/api/rainbow-table/%d", rainbowTable.ID))
		writer.WriteHeader(http.StatusNoContent)
	}
}
