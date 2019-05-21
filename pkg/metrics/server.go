package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func StartPrometheusMetricsServer(prometheusPort int) {
	go func() {
		http.Handle("/prometheus", promhttp.Handler())
		log.Infof("Starting prometheus server on port %d...", prometheusPort)
		err := http.ListenAndServe(fmt.Sprintf(":%d", prometheusPort), nil)

		if err != nil {
			log.Errorf("Failed to start prometheus server on port %d: %s", prometheusPort, err)
			os.Exit(1)
		}

	}()

}
