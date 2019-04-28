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
		writer.WriteHeader(400)
		writer.Write([]byte(fmt.Sprintf("Failed to parse integer query parameter %s: %s", parameter, value)))
		return parsedValue, err
	}

	return parsedValue, nil
}

func getIdPathParamValue(writer http.ResponseWriter, request *http.Request) (int16, error) {
	vars := mux.Vars(request)
	id, err := strconv.ParseInt(vars["id"], 10, 16)

	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(fmt.Sprintf("Failed to parse id from path param: %s", vars["id"])))
		return 0, err
	}

	return int16(id), nil
}
