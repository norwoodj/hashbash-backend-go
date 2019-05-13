package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"strconv"
)

func getIncludeNotFoundQueryParam(queryParameters url.Values) bool {
	return queryParameters.Get("includeNotFound") == "true"
}

func getIntQueryParamValue(
	queryParameters url.Values,
	parameter string,
	defaultValue int,
	writer http.ResponseWriter,
) (int, error) {
	value := queryParameters.Get(parameter)

	if value == "" {
		return defaultValue, nil
	}

	parsedValue, err := strconv.Atoi(value)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf("Failed to parse integer query parameter %s: %s", parameter, value)))
		return parsedValue, err
	}

	return parsedValue, nil
}

func getIdPathParamValue(
	idParamName string,
	writer http.ResponseWriter,
	request *http.Request,
	bitSize int,
) (interface{}, error) {
	vars := mux.Vars(request)
	id, err := strconv.ParseInt(vars[idParamName], 10, bitSize)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf("Failed to parse id from path param: %s", vars["id"])))
		return 0, err
	}

	return id, nil
}

func convertRainbowTableId(rainbowTableId interface{}) int16 {
	return int16(rainbowTableId.(int64))
}
