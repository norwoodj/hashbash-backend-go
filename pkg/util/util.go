package util

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
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

func SetupLogging() error {
	logLevel, err := log.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		return err
	}

	log.SetLevel(logLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	return nil
}

func GetTcpListenerOrDie(tcpAddr string) net.Listener {
	listener, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("Failed to listen on tcp address %s: %s", tcpAddr, err)
	}

	return listener
}

func GetUnixSocketListenerOrDie(socketPath string) net.Listener {
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("Failed to listen on unix socket %s: %s", socketPath, err)
	}

	return listener
}

func GetSystemdListenersOrDie(socketFdName string, listenersByName map[string][]net.Listener) []net.Listener {
	listener, ok := listenersByName[socketFdName]
	if !ok {
		log.Fatalf("No systemd socket found with fd name %s", socketFdName)
	}

	return listener
}

func WaitForSignalGracefulShutdown(cancel context.CancelFunc) {
	gracefulStop := make(chan os.Signal, 1)
    signal.Notify(gracefulStop, syscall.SIGTERM)
    signal.Notify(gracefulStop, syscall.SIGINT)

	shutdownSignal := <-gracefulStop
	log.Infof("Received signal %s, stopping servers...", shutdownSignal)
	cancel()
}

func GetServerForHandler(handler http.Handler) http.Server {
	readTimeout := viper.GetDuration("read-timeout")
	writeTimeout := viper.GetDuration("write-timeout")
	idleTimeout := viper.GetDuration("idle-timeout")

	return http.Server{
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      handler,
	}
}

func StartServer(server http.Server, listener net.Listener) error {
	log.Infof("Starting %s server...", listener.Addr().String())
	if err := server.Serve(listener); err != nil {
		log.Errorf("Error running %s server: %s", listener.Addr().String(), err)
		return err
	}

	return nil
}

func HandleServerShutdown(done context.Context, server http.Server, listener net.Listener) error {
	<-done.Done()

	serverName := listener.Addr().String()
	shutdownTimeout := viper.GetDuration("shutdown-timeout")
	log.Infof("Shutting down %s server gracefully with %s timeout", serverName, shutdownTimeout)
	ctx, _ := context.WithTimeout(context.Background(), shutdownTimeout)
	err := server.Shutdown(ctx)

	if err != nil {
		log.Errorf("Error shutting down %s server gracefully: %s", serverName, err)
		return err
	}

	log.Infof("Shut down %s server successfully", serverName)
	return nil
}

func GetManagementHandler() http.Handler {
	healthcheckHandler := healthcheck.NewHandler()
	prometheusHandler := promhttp.Handler()

	router := http.ServeMux{}
	router.Handle("/metrics", prometheusHandler)
	router.Handle("/", healthcheckHandler)
	return &router
}
