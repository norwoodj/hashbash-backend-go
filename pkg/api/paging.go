package api

import (
	"net/http"

	"github.com/norwoodj/hashbash-backend-go/pkg/database"
)

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
