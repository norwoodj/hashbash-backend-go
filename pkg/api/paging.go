package api

import (
	"github.com/norwoodj/hashbash-backend-go/pkg/service"
	"net/http"
)

func getPageConfigFromRequest(writer http.ResponseWriter, request *http.Request) (service.PageConfig, error) {
	queryParameters := request.URL.Query()

	pageNumber, err := getIntQueryParamValue(queryParameters, PagingQueryPageNumber, 0, writer)
	if err != nil {
		return service.PageConfig{}, err
	}

	pageSize, err := getIntQueryParamValue(queryParameters, PagingQueryPageSize, 10, writer)
	if err != nil {
		return service.PageConfig{}, err
	}

	sortKey := queryParameters.Get(PagingQuerySortKey)
	if sortKey == "" {
		sortKey = "id"
	}

	descending := queryParameters.Get(PagingQuerySortOrder) != PagingQuerySortOrderAscending

	return service.PageConfig{
		Descending: descending,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		SortKey:    sortKey,
	}, nil
}
