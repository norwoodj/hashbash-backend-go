package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var version string

func main() {
	rootCmd := newRootCommand(version)

	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
