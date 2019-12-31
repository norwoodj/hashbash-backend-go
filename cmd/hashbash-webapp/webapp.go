package main

import (
	"os"
	"strings"
	"sync"

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

func openOrCreate(logFilename string) (*os.File, error) {
	if _, err := os.Stat(logFilename); os.IsNotExist(err) {
		return os.Create(logFilename)
	}

	return os.Open(logFilename)
}

func setupLogging() {
	logLevel, _ := log.ParseLevel(viper.GetString("log-level"))
	log.SetLevel(logLevel)

	logFilename := viper.GetString("log-file")
	if logFilename != "" {
		file, err := openOrCreate(logFilename)

		if err != nil {
			log.Errorf("Could not open log file %s", logFilename)
		} else {
			log.SetOutput(file)
		}
	}
}

func hashbashWebapp(_ *cobra.Command, _ []string) {
	setupLogging()

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

	router := mux.NewRouter()
	api.AddRainbowTableRoutes(router, rainbowTableService, hashbashProducers)
	api.AddRainbowTableSearchRoutes(router, rainbowTableSearchService, hashbashProducers)

	frontendTemplatesDir := viper.GetString("frontend-template-path")
	err = frontend.RegisterTemplateHandler(router, frontendTemplatesDir)
	if err != nil {
		log.Errorf("Failed to read frontend directory %s: %s", frontendTemplatesDir, err)
		os.Exit(1)
	}

	walkRoutes(router)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(2)
	webPort := viper.GetInt("web-port")

	go util.StartManagementServer(&waitGroup)
	go util.StartHttpServer(webPort, "hashbash webapp", router, &waitGroup)
	waitGroup.Wait()
}
