package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

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

func WalkRoutes(router *mux.Router) {
	log.Debugf("Walking registered routes...")

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			log.Debugf("Route: %s", pathTemplate)
		}

		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			log.Debugf("Path regexp: %s", pathRegexp)
		}

		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			log.Debugf("Queries templates: [%s]", strings.Join(queriesTemplates, ","))
		}

		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			log.Debugf("Queries regexps: [%s]", strings.Join(queriesRegexps, ","))
		}

		methods, err := route.GetMethods()
		if err == nil {
			log.Debugf("Methods: [%s]", strings.Join(methods, ","))
		}

		return nil
	})
}
