package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/norwoodj/hashbash-backend-go/pkg/mq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func startConsumersAndHandleSignals(consumers mq.HashbashMqConsumers, shutdownGraceDuration time.Duration) {
	quit := make(chan bool)

	log.Infof("Starting hashbash consumers...")
	go consumers.DeleteRainbowTableConsumer.ConsumeMessages(quit)
	go consumers.GenerateRainbowTableConsumer.ConsumeMessages(quit)
	go consumers.SearchRainbowTableConsumer.ConsumeMessages(quit)

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

	connection, err := mq.AcquireMqConnection()
	defer connection.Close()

	if err != nil {
		log.Errorf("Failed to create rabbitmq connection: %s", err)
		os.Exit(1)
	}

	hashbashConsumers, err := mq.CreateConsumers(connection)
	if err != nil {
		log.Errorf("Failed to instantiate rabbitmq consumers: %s", err)
		os.Exit(1)
	}

	startConsumersAndHandleSignals(hashbashConsumers, viper.GetDuration("shutdown-timeout"))
}
