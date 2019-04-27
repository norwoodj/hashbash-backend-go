package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/norwoodj/hashbash-backend-go/pkg/database"
)

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

func getPageConfigFromRequest(writer http.ResponseWriter, request *http.Request) (database.PageConfig, error) {
	queryParameters := request.URL.Query()

	pageNumber, err := getIntQueryParamValue(queryParameters, "pageNumber", 0, writer)
	if err != nil {
		return database.PageConfig{}, err
	}

	pageSize, err := getIntQueryParamValue(queryParameters, "pageSize", 10, writer)
	if err != nil {
		return database.PageConfig{}, err
	}

	sortKey := queryParameters.Get("sortKey")
	if sortKey == "" {
		sortKey = "id"
	}

	descending := queryParameters.Get("sortOrder") != "ASC"

	return database.PageConfig{
		Descending: descending,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		SortKey:    sortKey,
	}, nil
}
