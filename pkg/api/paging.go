package api

import (
	"net/http"

	"github.com/norwoodj/hashbash-backend-go/pkg/database"
	"github.com/norwoodj/hashbash-backend-go/pkg/util"
)

func getPageConfigFromRequest(writer http.ResponseWriter, request *http.Request) (database.PageConfig, error) {
	queryParameters := request.URL.Query()

	pageNumber, err := getIntQueryParamValue(queryParameters, util.PagingQueryPageNumber, 0, writer)
	if err != nil {
		return database.PageConfig{}, err
	}

	pageSize, err := getIntQueryParamValue(queryParameters, util.PagingQueryPageSize, 10, writer)
	if err != nil {
		return database.PageConfig{}, err
	}

	sortKey := queryParameters.Get(util.PagingQuerySortKey)
	if sortKey == "" {
		sortKey = "id"
	}

	descending := queryParameters.Get(util.PagingQuerySortOrder) != util.PagingQuerySortOrderAscending

	return database.PageConfig{
		Descending: descending,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		SortKey:    sortKey,
	}, nil
}
