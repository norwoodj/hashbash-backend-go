package main

import (
	"strings"
	"time"

	"github.com/norwoodj/hashbash-backend-go/pkg/dao"
	"github.com/norwoodj/hashbash-backend-go/pkg/frontend"
	"github.com/norwoodj/hashbash-backend-go/pkg/rabbitmq"
	"github.com/norwoodj/hashbash-backend-go/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newRootCommand(version string) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:     "hashbash-webapp",
		Short:   "Serve the hashbash web application",
		Run:     hashbashWebapp,
		Version: version,
	}

	webappFlags := rootCmd.PersistentFlags()
	webappFlags.Duration("read-timeout", time.Second*5, "The duration for which the server waits for read operations to complete, e.g. 15s or 1m")
	webappFlags.Duration("write-timeout", time.Second*5, "The duration for which the server waits for write operations to complete, e.g. 15s or 1m")
	webappFlags.Duration("idle-timeout", 0, "The duration for which the server waits for idle connections to write/read data")
	webappFlags.Duration("shutdown-timeout", time.Second*10, "The duration for which the server waits for existing connections to complete at shutdown time")
	webappFlags.StringSlice("http-addr", []string{}, "host:port interfaces on which to serve http traffic, may be repeated")
	webappFlags.StringSlice("http-name", []string{}, "systemd socket name on which to serve http traffic, may be repeated")
	webappFlags.StringSlice("http-sock", []string{}, "File paths of sockets on which to serve http traffic, may be repeated")

	util.AddDefaultFlags(webappFlags)
	dao.AddDatabaseFlags(webappFlags)
	frontend.AddFrontendFlags(webappFlags)
	rabbitmq.AddRabbitMqFlags(webappFlags)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("HASHBASH")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	_ = viper.BindPFlags(rootCmd.PersistentFlags())

	return rootCmd
}
