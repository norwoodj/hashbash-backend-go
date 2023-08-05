package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var buildTimestamp string
var gitRevision string
var version string

func main() {
	rootCmd := newRootCommand(buildTimestamp, gitRevision, version)

	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
