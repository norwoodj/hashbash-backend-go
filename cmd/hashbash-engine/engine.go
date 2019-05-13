package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/norwoodj/hashbash-backend-go/pkg/rabbitmq"
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

	//db := database.GetConnectionOrDie()
	//rainbowTableService := service.NewRainbowTableService(db)
	//rainbowTableSearchService := service.NewRainbowTableSearchService(db)

	connection := rabbitmq.AcquireMqConnectionOrDie()
	defer connection.Close()

	hashbashConsumers, err := rabbitmq.CreateConsumerWorkers(connection)
	if err != nil {
		log.Errorf("Failed to instantiate rabbitmq consumers: %s", err)
		os.Exit(1)
	}

	startConsumersAndHandleSignals(hashbashConsumers, viper.GetDuration("shutdown-timeout"))
}