package main

import (
	"context"
	"github.com/coreos/go-systemd/activation"
	"github.com/gorilla/handlers"
	"golang.org/x/sync/errgroup"
	"os"

	"github.com/norwoodj/hashbash-backend-go/pkg/dao"
	"github.com/norwoodj/hashbash-backend-go/pkg/metrics"
	"github.com/norwoodj/hashbash-backend-go/pkg/rabbitmq"
	"github.com/norwoodj/hashbash-backend-go/pkg/rainbow"
	"github.com/norwoodj/hashbash-backend-go/pkg/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func toErrFunc(f func(chan error)) func() error {
	startErrors := make(chan error)

	return func() error {
		f(startErrors)

		for e := range startErrors {
			return e
		}

		return nil
	}
}

func startConsumers(
	consumers rabbitmq.HashbashMqConsumerWorkers,
	startErrGroup *errgroup.Group,
	shutdownErrGroup *errgroup.Group,
) chan bool {
	quit := make(chan bool)

	log.Info().Msg("Starting hashbash consumers...")
	startErrGroup.Go(toErrFunc(
		func(startErrors chan error) {
			consumers.HashbashDeleteRainbowTableConsumer.ConsumeMessages(quit, startErrors)
		},
	))

	startErrGroup.Go(toErrFunc(
		func(startErrors chan error) {
			consumers.HashbashGenerateRainbowTableConsumer.ConsumeMessages(quit, startErrors)
		},
	))

	startErrGroup.Go(toErrFunc(
		func(startErrors chan error) {
			consumers.HashbashSearchRainbowTableConsumer.ConsumeMessages(quit, startErrors)
		},
	))

	return quit
}

func registerShutdownConsumers(quit chan bool, ctx context.Context) {
	<-ctx.Done()
	close(quit)
}

func hashbashEngine(_ *cobra.Command, _ []string) {
	err := util.SetupLogging()
	if err != nil {
		log.Error().Err(err).Msg("Failed to setup logging")
		os.Exit(1)
	}

	dbEngine := viper.GetString("database-engine")
	db := dao.GetConnectionOrDie(dbEngine)
	rainbowTableService := dao.NewRainbowTableService(db)
	rainbowChainService := dao.NewRainbowChainService(db, dbEngine)
	rainbowTableSearchService := dao.NewRainbowTableSearchService(db)

	generateJobConfig := rainbow.TableGenerateJobConfig{
		ChainBatchSize: viper.GetInt64("generate-batch-size"),
		NumThreads:     viper.GetInt("generate-num-threads"),
	}

	searchJobConfig := rainbow.TableSearchJobConfig{
		SearchHashBatchSize: viper.GetInt("search-batch-size"),
		NumThreads:          viper.GetInt("search-num-threads"),
	}

	chainGenerationSummary := metrics.NewRainbowChainSummary("chain", "generate_seconds")
	chainWriteSummary := metrics.NewRainbowChainSummary("chain", "write_seconds")
	chainsCreatedCounter := metrics.NewRainbowChainCounter("chain", "created_total")

	rainbowTableGenerateJobService := rainbow.NewRainbowTableGeneratorJobService(
		generateJobConfig,
		rainbowChainService,
		rainbowTableService,
		chainGenerationSummary,
		chainWriteSummary,
		chainsCreatedCounter,
	)

	rainbowTableSearchJobService := rainbow.NewRainbowTableSearchJobService(
		searchJobConfig,
		rainbowChainService,
		rainbowTableService,
		rainbowTableSearchService,
	)

	connection := rabbitmq.AcquireMqConnectionOrDie()
	defer connection.Close()

	hashbashConsumers, err := rabbitmq.CreateConsumerWorkers(
		connection,
		rainbowTableService,
		rainbowTableGenerateJobService,
		rainbowTableSearchJobService,
	)

	if err != nil {
		log.Error().Err(err).Msg("Failed to instantiate rabbitmq consumers")
		os.Exit(1)
	}

	done, cancel := context.WithCancel(context.Background())
	startErrGroup, _ := errgroup.WithContext(done)
	shutdownErrGroup, _ := errgroup.WithContext(done)

	quit := startConsumers(hashbashConsumers, startErrGroup, shutdownErrGroup)
	go registerShutdownConsumers(quit, done)

	systemdListenersByName, err := activation.ListenersWithNames()

	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to retrieve systemd sockets by name")
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
