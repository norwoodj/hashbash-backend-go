package main

import (
	"context"
	"github.com/coreos/go-systemd/activation"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/norwoodj/hashbash-backend-go/pkg/api"
	"github.com/norwoodj/hashbash-backend-go/pkg/dao"
	"github.com/norwoodj/hashbash-backend-go/pkg/frontend"
	"github.com/norwoodj/hashbash-backend-go/pkg/rabbitmq"
	"github.com/norwoodj/hashbash-backend-go/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

func walkRoutes(router *mux.Router) {
	log.Debugf("Walking registered routes...")

	_ = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
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

func startHttpHandler(startErrGroup *errgroup.Group, shutdownErrGroup *errgroup.Group, done context.Context, listener net.Listener, handler http.Handler) {
	server := util.GetServerForHandler(handler)
	startErrGroup.Go(func() error { return util.StartServer(server, listener) })
	shutdownErrGroup.Go(func() error { return util.HandleServerShutdown(done, server, listener) })
}

func hashbashWebapp(_ *cobra.Command, _ []string) {
	err := util.SetupLogging()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	util.DoInitialDelay()

	db := dao.GetConnectionOrDie(viper.GetString("database-engine"))
	rainbowTableService := dao.NewRainbowTableService(db)
	rainbowTableSearchService := dao.NewRainbowTableSearchService(db)

	connection := rabbitmq.AcquireMqConnectionOrDie()
	defer connection.Close()

	hashbashProducers, err := rabbitmq.CreateProducers(connection)
	if err != nil {
		log.Fatalf("Failed to instantiate rabbitmq producers: %s", err)
	}

	router := mux.NewRouter()
	api.AddRainbowTableRoutes(router, rainbowTableService, hashbashProducers)
	api.AddRainbowTableSearchRoutes(router, rainbowTableSearchService, hashbashProducers)

	frontendTemplatesDir := viper.GetString("frontend-template-path")
	err = frontend.RegisterTemplateHandler(router, frontendTemplatesDir)
	if err != nil {
		log.Fatalf("Failed to read frontend directory %s: %s", frontendTemplatesDir, err)
	}

	walkRoutes(router)
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	systemdListenersByName, err := activation.ListenersWithNames()

	if err != nil {
		log.Fatalf("Failed to retrieve systemd sockets by name: %s", err)
	}

	done, cancel := context.WithCancel(context.Background())
	startErrGroup, _ := errgroup.WithContext(done)
	shutdownErrGroup, _ := errgroup.WithContext(done)

	for _, addr := range viper.GetStringSlice("http-addr") {
		listener := util.GetTcpListenerOrDie(addr)
		startHttpHandler(startErrGroup, shutdownErrGroup, done, listener, loggedRouter)
	}

	for _, socketPath := range viper.GetStringSlice("http-sock") {
		listener := util.GetUnixSocketListenerOrDie(socketPath)
		startHttpHandler(startErrGroup, shutdownErrGroup, done, listener, loggedRouter)
	}

	for _, socketFdName := range viper.GetStringSlice("http-name") {
		listeners := util.GetSystemdListenersOrDie(socketFdName, systemdListenersByName)

		for _, l := range listeners {
			startHttpHandler(startErrGroup, shutdownErrGroup, done, l, loggedRouter)
		}
	}

	managementHandler := handlers.LoggingHandler(os.Stdout, util.GetManagementHandler())
	for _, addr := range viper.GetStringSlice("management-addr") {
		listener := util.GetTcpListenerOrDie(addr)
		startHttpHandler(startErrGroup, shutdownErrGroup, done, listener, managementHandler)
	}

	for _, socketPath := range viper.GetStringSlice("management-sock") {
		listener := util.GetUnixSocketListenerOrDie(socketPath)
		startHttpHandler(startErrGroup, shutdownErrGroup, done, listener, managementHandler)
	}

	for _, socketFdName := range viper.GetStringSlice("management-name") {
		listeners := util.GetSystemdListenersOrDie(socketFdName, systemdListenersByName)

		for _, l := range listeners {
			startHttpHandler(startErrGroup, shutdownErrGroup, done, l, managementHandler)
		}
	}

	go util.WaitForSignalGracefulShutdown(cancel)

	go func() {
		if err := startErrGroup.Wait(); err != nil {
			log.Fatalf("Failed to start servers: %s", err)
		}
	}()

	if err := shutdownErrGroup.Wait(); err != nil {
		log.Fatalf("Error shutting down servers: %s", err)
	}

	log.Info("Shutdown successful")
}
