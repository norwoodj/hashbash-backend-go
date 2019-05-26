package metrics

import (
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
)

func GetRainbowTableMetricLabels(rainbowTable *model.RainbowTable, batchSize int) prometheus.Labels {
	return prometheus.Labels{
		"batch_size":       strconv.Itoa(batchSize),
		"chain_length":     strconv.Itoa(int(rainbowTable.ChainLength)),
		"hash_function":    rainbowTable.HashFunction,
		"rainbow_table_id": strconv.Itoa(int(rainbowTable.ID)),
	}
}

func NewRainbowChainSummary(
	metricSubsystem string,
	metricName string,
) *prometheus.SummaryVec {
	summaryOpts := prometheus.SummaryOpts{
		Namespace:   "rainbow",
		Subsystem:   metricSubsystem,
		Name:        metricName,
	}

	return promauto.NewSummaryVec(summaryOpts, []string{
		"batch_size",
		"chain_length",
		"hash_function",
		"rainbow_table_id",
	})
}

func NewRainbowChainCounter(
	metricSubsystem string,
	metricName string,
) *prometheus.CounterVec {
	counterOpts := prometheus.CounterOpts{
		Namespace:   "rainbow",
		Subsystem:   metricSubsystem,
		Name:        metricName,
	}

	return promauto.NewCounterVec(counterOpts, []string{
		"batch_size",
		"chain_length",
		"hash_function",
		"rainbow_table_id",
	})
}
