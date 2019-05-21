package util

import (
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
	"time"
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
