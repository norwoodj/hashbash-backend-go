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
	flags.StringSlice("management-addr", []string{}, "host:port interfaces on which to serve management http traffic, may be repeated")
	flags.StringSlice("management-name", []string{}, "systemd socket name on which to serve management http traffic, may be repeated")
	flags.StringSlice("management-sock", []string{}, "File paths of sockets on which to serve management http traffic, may be repeated")
}
