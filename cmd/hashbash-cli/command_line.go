package main

import (
	"strings"

	"github.com/norwoodj/hashbash-backend-go/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newRootCommand(version string) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:     "hashbash-cli",
		Short:   "CLI for automating various hashbash related tasks",
		Version: version,
	}

	util.AddDefaultFlags(rootCmd.PersistentFlags())
	rootCmd.AddCommand(newSearchSubcommand())

	viper.AutomaticEnv()
	viper.SetEnvPrefix("HASHBASH")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.BindPFlags(rootCmd.PersistentFlags())
	return rootCmd
}
