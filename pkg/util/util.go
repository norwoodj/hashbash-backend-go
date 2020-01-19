package util

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/heptiolabs/healthcheck"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func IntOrDefault(value int64, defaultValue int64) int64 {
	if value == 0 {
		return defaultValue
	}

	return value
}

func StringOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	return value
}

func DoInitialDelay() {
	initialDelay := viper.GetDuration("initial-delay")

	if initialDelay != 0 {
		log.Infof("Delaying startup by %s, to allow for mysql/rabbitmq to start up...", initialDelay)
		time.Sleep(initialDelay)
	}
}

func SetupLogging() (*os.File, bool, error) {
	logLevel, err := log.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse log level %s: %s", logLevel, err)
	}

	log.SetLevel(logLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	logFilename := viper.GetString("log-file")
	if logFilename == "" {
		return os.Stdout, false, nil
	}

	file, err := os.OpenFile(logFilename, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
	if err != nil {
		return nil, false, fmt.Errorf("could not open log file %s: %s", logFilename, err)
	}

	log.SetOutput(file)
	return file, true, nil
}

func startServerAndHandleSignals(server *http.Server, serverName string, port int, shutdownTimeout time.Duration) {
	go func() {
		log.Infof("Starting %s server on port %d...", serverName, port)
		if err := server.ListenAndServe(); err != nil {
			log.Errorf("Error running %s server: %s", serverName, err)
			os.Exit(1)
		}
	}()

	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	shutdownSignal := <-gracefulStop

	log.Infof("Received Signal %s, shutting down %s server gracefully with %s timeout", shutdownSignal, serverName, shutdownTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		log.Errorf("Error shutting down %s server gracefully: %s", serverName, err)
		os.Exit(1)
	}

	os.Exit(0)
}

func StartManagementServer(waitGroup *sync.WaitGroup) {
	managementPort := viper.GetInt("management-port")
	healthcheckHandler := healthcheck.NewHandler()
	prometheusHandler := promhttp.Handler()

	router := http.ServeMux{}
	router.Handle("/prometheus", prometheusHandler)
	router.Handle("/", healthcheckHandler)

	StartHttpServer(managementPort, "management", &router, waitGroup)
}

func StartHttpServer(port int, serverName string, handler http.Handler, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	startServerAndHandleSignals(server, serverName, port, time.Second*5)
}
