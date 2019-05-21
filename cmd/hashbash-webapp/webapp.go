package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/norwoodj/hashbash-backend-go/pkg/api"
	"github.com/norwoodj/hashbash-backend-go/pkg/dao"
	"github.com/norwoodj/hashbash-backend-go/pkg/database"
	"github.com/norwoodj/hashbash-backend-go/pkg/frontend"
	"github.com/norwoodj/hashbash-backend-go/pkg/rabbitmq"
	"github.com/norwoodj/hashbash-backend-go/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
)

func startServerAndHandleSignals(server *http.Server, port int, shutdownTimeout time.Duration) {
	go func() {
		log.Infof("Starting hashbash webapp on port %d...", port)
		if err := server.ListenAndServe(); err != nil {
			log.Errorf("Error running hashbash webapp: %s", err)
			os.Exit(1)
		}
	}()

	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	shutdownSignal := <-gracefulStop

	log.Infof("Received Signal %s, shutting down gracefully with %s timeout", shutdownSignal, shutdownTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		log.Errorf("Error shutting down gracefully, some connections may have been dropped")
		os.Exit(1)
	}

	os.Exit(0)
}

func walkRoutes(router *mux.Router) {
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

func hashbashWebapp(_ *cobra.Command, _ []string) {
	logLevel, _ := log.ParseLevel(viper.GetString("log-level"))
	log.SetLevel(logLevel)

	util.DoInitialDelay()
	db := database.GetConnectionOrDie()
	rainbowTableService := dao.NewRainbowTableService(db)
	rainbowTableSearchService := dao.NewRainbowTableSearchService(db)

	connection := rabbitmq.AcquireMqConnectionOrDie()
	defer connection.Close()

	hashbashProducers, err := rabbitmq.CreateProducers(connection)
	if err != nil {
		log.Errorf("Failed to instantiate rabbitmq producers: %s", err)
		os.Exit(1)
	}

	port := viper.GetInt("web-port")
	router := mux.NewRouter()
	server := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	api.AddRainbowTableRoutes(router, rainbowTableService, hashbashProducers)
	api.AddRainbowTableSearchRoutes(router, rainbowTableSearchService, hashbashProducers)

	frontendTemplatesDir := viper.GetString("frontend-template-path")
	err = frontend.RegisterTemplateHandler(router, frontendTemplatesDir)
	if err != nil {
		log.Errorf("Failed to read frontend directory %s: %s", frontendTemplatesDir, err)
		os.Exit(1)
	}

	walkRoutes(router)
	startServerAndHandleSignals(&server, port, viper.GetDuration("shutdown-timeout"))
}
