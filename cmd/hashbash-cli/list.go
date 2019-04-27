package main

import (
	"os"
	"strconv"

	"github.com/norwoodj/hashbash-backend-go/pkg/database"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func rainbowTableToRow(rainbowTable *model.RainbowTable) []string {
	return []string{
		rainbowTable.Name,
		rainbowTable.Status,
		strconv.FormatInt(rainbowTable.NumChains, 10),
		strconv.FormatInt(rainbowTable.ChainsGenerated, 10),
		strconv.FormatInt(rainbowTable.ChainLength, 10),
		rainbowTable.HashFunction,
		strconv.FormatInt(rainbowTable.PasswordLength, 10),
		rainbowTable.CharacterSet,
	}
}

func printRainbowTableTable(rainbowTables []model.RainbowTable) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Name",
		"Status",
		"Num Chains",
		"Chains Generated",
		"Chain Length",
		"Hash Function",
		"Password Length",
		"Character Set",
	})

	for _, r := range rainbowTables {
		table.Append(rainbowTableToRow(&r))
	}

	table.Render()
}

func listSubcommand(_ *cobra.Command, _ []string) {
	db := database.GetConnectionOrDie()
	defer db.Close()

	rainbowTables := make([]model.RainbowTable, 0)
	db.Find(&rainbowTables)
	printRainbowTableTable(rainbowTables)
}
