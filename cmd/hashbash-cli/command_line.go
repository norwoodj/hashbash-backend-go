package main

import (
	"strings"

	"github.com/norwoodj/hashbash-backend-go/pkg/database"
	"github.com/norwoodj/hashbash-backend-go/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newRootCommand() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "hashbash-cli",
		Short: "CLI for automating various hashbash related tasks",
	}

	util.AddDefaultFlags(rootCmd.PersistentFlags())
	database.AddDatabaseFlags(rootCmd.PersistentFlags())

	rootCmd.AddCommand(newListSubcommand())

	viper.AutomaticEnv()
	viper.SetEnvPrefix("HASHBASH")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.BindPFlags(rootCmd.PersistentFlags())
	return rootCmd
}
