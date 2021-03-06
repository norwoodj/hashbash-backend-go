package util

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func possibleLogLevels() []string {
	levels := make([]string, 0)

	for _, l := range log.AllLevels {
		levels = append(levels, l.String())
	}

	return levels
}

func AddDefaultFlags(flags *pflag.FlagSet) {
	logLevelUsage := fmt.Sprintf("Level of logs that should printed, one of (%s)", strings.Join(possibleLogLevels(), ", "))
	flags.StringP("log-level", "l", "info", logLevelUsage)
	flags.Duration("initial-delay", 0, "Time to delay startup to allow for database/rabbit to start (for local docker dev)")
	flags.IntP("management-port", "m", 8081, "The port on which to serve the management server")
}
