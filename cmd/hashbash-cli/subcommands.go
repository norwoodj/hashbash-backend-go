package main

import "github.com/spf13/cobra"

func newListSubcommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List current rainbow tables from the database",
		Run:   listSubcommand,
	}
}
