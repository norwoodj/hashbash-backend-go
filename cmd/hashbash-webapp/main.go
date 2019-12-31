package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var version string

func main() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	rootCmd := newRootCommand(version)

	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
