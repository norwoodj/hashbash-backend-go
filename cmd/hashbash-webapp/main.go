package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var Version = "UNVERSIONED"

func main() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	rootCmd := newRootCommand(Version)

	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
