package main

import (
	"context"
	"github.com/coreos/go-systemd/activation"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/norwoodj/hashbash-backend-go/pkg/api"
	"github.com/norwoodj/hashbash-backend-go/pkg/dao"
	"github.com/norwoodj/hashbash-backend-go/pkg/frontend"
	"github.com/norwoodj/hashbash-backend-go/pkg/rabbitmq"
	"github.com/norwoodj/hashbash-backend-go/pkg/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

func walkRoutes(router *mux.Router) {
	log.Debug().Msgf("Walking registered routes...")

	_ = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			log.Debug().Msgf("Route: %s", pathTemplate)
		}

		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			log.Debug().Msgf("Path regexp: %s", pathRegexp)
		}

		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			log.Debug().Msgf("Queries templates: [%s]", strings.Join(queriesTemplates, ","))
		}

		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			log.Debug().Msgf("Queries regexps: [%s]", strings.Join(queriesRegexps, ","))
		}

		methods, err := route.GetMethods()
		if err == nil {
			log.Debug().Msgf("Methods: [%s]", strings.Join(methods, ","))
		}

		return nil
	})
}

func hashbashWebapp(buildTimestamp string, gitRevision string, version string) {
	err := util.SetupLogging()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to setup logging")
	}

	db := dao.GetConnectionOrDie(viper.GetString("database-engine"))
	rainbowTableService := dao.NewRainbowTableService(db)
	rainbowTableSearchService := dao.NewRainbowTableSearchService(db)

	connection := rabbitmq.AcquireMqConnectionOrDie()
	defer connection.Close()

	hashbashProducers, err := rabbitmq.CreateProducers(connection)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to instantiate rabbitmq producers")
	}

	router := mux.NewRouter()
	api.AddRainbowTableRoutes(router, rainbowTableService, hashbashProducers)
	api.AddRainbowTableSearchRoutes(router, rainbowTableSearchService, hashbashProducers)
	api.AddVersionRoutes(router, buildTimestamp, gitRevision, version)

	frontendTemplatesDir := viper.GetString("frontend-template-path")
	err = frontend.RegisterTemplateHandler(router, frontendTemplatesDir)
	if err != nil {
		log.Fatal().
			Err(err).
			Msgf("Failed to read frontend directory %s", frontendTemplatesDir)
	}

	walkRoutes(router)
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	systemdListenersByName, err := activation.ListenersWithNames()

	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to retrieve systemd sockets by name")
	}

	done, cancel := context.WithCancel(context.Background())
	startErrGroup, _ := errgroup.WithContext(done)
	shutdownErrGroup, _ := errgroup.WithContext(done)

	for _, addr := range viper.GetStringSlice("http-addr") {
		listener := util.GetTcpListenerOrDie(addr)
		util.StartHttpHandler(startErrGroup, shutdownErrGroup, done, listener, loggedRouter)
	}

	for _, socketPath := range viper.GetStringSlice("http-sock") {
		listener := util.GetUnixSocketListenerOrDie(socketPath)
		util.StartHttpHandler(startErrGroup, shutdownErrGroup, done, listener, loggedRouter)
	}

	for _, socketFdName := range viper.GetStringSlice("http-name") {
		listeners := util.GetSystemdListenersOrDie(socketFdName, systemdListenersByName)

		for _, l := range listeners {
			util.StartHttpHandler(startErrGroup, shutdownErrGroup, done, l, loggedRouter)
		}
	}

	managementHandler := handlers.LoggingHandler(os.Stdout, util.GetManagementHandler())
	for _, addr := range viper.GetStringSlice("management-addr") {
		listener := util.GetTcpListenerOrDie(addr)
		util.StartHttpHandler(startErrGroup, shutdownErrGroup, done, listener, managementHandler)
	}

	for _, socketPath := range viper.GetStringSlice("management-sock") {
		listener := util.GetUnixSocketListenerOrDie(socketPath)
		util.StartHttpHandler(startErrGroup, shutdownErrGroup, done, listener, managementHandler)
	}

	for _, socketFdName := range viper.GetStringSlice("management-name") {
		listeners := util.GetSystemdListenersOrDie(socketFdName, systemdListenersByName)

		for _, l := range listeners {
			util.StartHttpHandler(startErrGroup, shutdownErrGroup, done, l, managementHandler)
		}
	}

	util.WaitForSignalGracefulShutdown(cancel, startErrGroup, shutdownErrGroup)
}
