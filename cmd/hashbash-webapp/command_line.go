package main

import (
	"strings"
	"time"

	"github.com/norwoodj/hashbash-backend-go/pkg/database"
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
	webappFlags.DurationP("shutdown-timeout", "s", time.Second*15, "The duration for which the server waits for existing connections to finish, e.g. 15s or 1m")
	webappFlags.IntP("web-port", "w", 8080, "Port on which to serve the hashbash webapp")

	util.AddDefaultFlags(webappFlags)
	database.AddDatabaseFlags(webappFlags)
	frontend.AddFrontendFlags(webappFlags)
	rabbitmq.AddRabbitMqFlags(webappFlags)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("HASHBASH")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.BindPFlags(rootCmd.PersistentFlags())
	return rootCmd
}
