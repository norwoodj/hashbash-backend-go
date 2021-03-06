package main

import (
	"github.com/norwoodj/hashbash-backend-go/pkg/dao"
	"github.com/norwoodj/hashbash-backend-go/pkg/rabbitmq"
	"github.com/norwoodj/hashbash-backend-go/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func newRootCommand(version string) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:     "hashbash-engine",
		Short:   "Run the hashbash consumers for servicing rainbow table requests",
		Run:     hashbashEngine,
		Version: version,
	}

	engineFlags := rootCmd.PersistentFlags()
	engineFlags.DurationP("shutdown-timeout", "s", time.Second*3, "The duration for which the metrics server waits for existing connections to finish, e.g. 15s or 1m")
	engineFlags.Int64("generate-batch-size", 1000, "The size of rainbow chain batches to generate at a time")
	engineFlags.Int("generate-num-threads", 8, "The number of threads to use when generating rainbow tables")
	engineFlags.Int("search-batch-size", 1000, "The size of rainbow chain batches to use in each search query")
	engineFlags.Int("search-num-threads", 8, "The number of threads to use when searching rainbow tables")

	util.AddDefaultFlags(engineFlags)
	dao.AddDatabaseFlags(engineFlags)
	rabbitmq.AddRabbitMqFlags(engineFlags)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("HASHBASH")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.BindPFlags(engineFlags)
	return rootCmd
}
