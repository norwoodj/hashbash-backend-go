package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/norwoodj/hashbash-backend-go/pkg/api"
	"github.com/norwoodj/hashbash-backend-go/pkg/database"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func hashbashWebapp(_ *cobra.Command, _ []string) {
	db := database.GetConnectionOrDie()
	port := viper.GetInt("web-port")
	router := mux.NewRouter()

	server := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	api.AddRainbowTableRoutes(router, db)
	startServerAndHandleSignals(&server, port, viper.GetDuration("shutdown-timeout"))
}
