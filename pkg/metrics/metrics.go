package metrics

import (
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
)

func GetRainbowTableMetricLabels(rainbowTable *model.RainbowTable) prometheus.Labels {
	return prometheus.Labels{
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
		ConstLabels: prometheus.Labels{"language": "golang"},
	}

	return promauto.NewSummaryVec(summaryOpts, []string{
		"chain_length",
		"hash_function",
		"rainbow_table_id",
	})
}
