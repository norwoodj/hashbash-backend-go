package api

import (
	"github.com/norwoodj/hashbash-backend-go/pkg/dao"
	"net/http"
)

func getPageConfigFromRequest(writer http.ResponseWriter, request *http.Request) (dao.PageConfig, error) {
	queryParameters := request.URL.Query()

	pageNumber, err := getIntQueryParamValue(queryParameters, PagingQueryPageNumber, 0, writer)
	if err != nil {
		return dao.PageConfig{}, err
	}

	pageSize, err := getIntQueryParamValue(queryParameters, PagingQueryPageSize, 10, writer)
	if err != nil {
		return dao.PageConfig{}, err
	}

	sortKey := queryParameters.Get(PagingQuerySortKey)
	if sortKey == "" {
		sortKey = "id"
	}

	descending := queryParameters.Get(PagingQuerySortOrder) != PagingQuerySortOrderAscending

	return dao.PageConfig{
		Descending: descending,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		SortKey:    sortKey,
	}, nil
}
