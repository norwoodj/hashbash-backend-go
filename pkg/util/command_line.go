package util

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
)

func possibleLogLevels() []string {
	levels := []zerolog.Level{
		zerolog.DebugLevel,
		zerolog.InfoLevel,
		zerolog.WarnLevel,
		zerolog.ErrorLevel,
		zerolog.FatalLevel,
		zerolog.PanicLevel,
		zerolog.NoLevel,
		zerolog.Disabled,
		zerolog.TraceLevel,
	}

	possibleCliValues := make([]string, 0)

	for _, l := range levels {
		possibleCliValues = append(possibleCliValues, l.String())
	}

	return possibleCliValues
}

func AddDefaultFlags(flags *pflag.FlagSet) {
	logLevelUsage := fmt.Sprintf("Level of logs that should printed, one of (%s)", strings.Join(possibleLogLevels(), ", "))
	flags.StringP("log-level", "l", "info", logLevelUsage)
	flags.StringSlice("management-addr", []string{}, "host:port interfaces on which to serve management http traffic, may be repeated")
	flags.StringSlice("management-name", []string{}, "systemd socket name on which to serve management http traffic, may be repeated")
	flags.StringSlice("management-sock", []string{}, "File paths of sockets on which to serve management http traffic, may be repeated")
}
