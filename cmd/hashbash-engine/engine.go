package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/norwoodj/hashbash-backend-go/pkg/dao"
	"github.com/norwoodj/hashbash-backend-go/pkg/database"
	"github.com/norwoodj/hashbash-backend-go/pkg/metrics"
	"github.com/norwoodj/hashbash-backend-go/pkg/rabbitmq"
	"github.com/norwoodj/hashbash-backend-go/pkg/rainbow"
	"github.com/norwoodj/hashbash-backend-go/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func startConsumersAndHandleSignals(consumers rabbitmq.HashbashMqConsumerWorkers, shutdownGraceDuration time.Duration) {
	consumerStartErrorChannels := []chan error{make(chan error), make(chan error), make(chan error)}
	quit := make(chan bool)

	log.Infof("Starting hashbash consumers...")
	go consumers.HashbashDeleteRainbowTableConsumer.ConsumeMessages(quit, consumerStartErrorChannels[0])
	go consumers.HashbashGenerateRainbowTableConsumer.ConsumeMessages(quit, consumerStartErrorChannels[1])
	go consumers.HashbashSearchRainbowTableConsumer.ConsumeMessages(quit, consumerStartErrorChannels[2])

	for _, errorChannel := range consumerStartErrorChannels {
		consumerStartError := <-errorChannel

		if consumerStartError != nil {
			log.Error(consumerStartError)
			os.Exit(1)
		}
	}

	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	shutdownSignal := <-gracefulStop

	log.Infof("Received Signal %s, shutting down gracefully, waiting %s for channels to close...", shutdownSignal, shutdownGraceDuration)
	close(quit)
	time.Sleep(shutdownGraceDuration)
}

func hashbashEngine(_ *cobra.Command, _ []string) {
	logLevel, _ := log.ParseLevel(viper.GetString("log-level"))
	log.SetLevel(logLevel)

	util.DoInitialDelay()

	db := database.GetConnectionOrDie()
	rainbowTableService := dao.NewRainbowTableService(db)
	rainbowChainService := dao.NewRainbowChainService(db)
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
	prometheusPort := viper.GetInt("prometheus-port")
	metrics.StartPrometheusMetricsServer(prometheusPort)

	rainbowTableGenerateJobService := rainbow.NewRainbowTableGeneratorJobService(
		generateJobConfig,
		rainbowChainService,
		rainbowTableService,
		chainGenerationSummary,
		chainWriteSummary,
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
		log.Errorf("Failed to instantiate rabbitmq consumers: %s", err)
		os.Exit(1)
	}

	startConsumersAndHandleSignals(hashbashConsumers, viper.GetDuration("shutdown-timeout"))
}
