package main

import (
	"encoding/json"
	"fmt"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
	"github.com/norwoodj/hashbash-backend-go/pkg/rainbow"
	"github.com/norwoodj/hashbash-backend-go/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func searchCommandFn(_ *cobra.Command, args []string) {
	rainbowTableId, err := strconv.ParseInt(args[0], 10, 16)
	if err != nil {
		log.Errorf("Failed to parse rainbow table ID %s as int: %s", args[0], err)
		os.Exit(1)
	}

	numPasswords, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		log.Errorf("Failed to parse number of passwords %s as int: %s", args[1], err)
		os.Exit(1)
	}

	hashbashHost := viper.GetString("hashbash-host")
	resp, err := http.Get(fmt.Sprintf("%s/api/rainbow-table/%d", hashbashHost, rainbowTableId))

	if err != nil {
		log.Errorf("Failed to retrieve rainbow table with ID %d", rainbowTableId)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Failed to read rainbow table response body: %s", err)
		os.Exit(1)
	}

	var rainbowTable model.RainbowTable
	err = json.Unmarshal(body, &rainbowTable)

	if err != nil {
		log.Errorf("Failed to parse json rainbow table response body: %s", err)
		os.Exit(1)
	}

	hashFunctionProvider, err := rainbow.GetHashFunctionProvider(rainbowTable.HashFunction)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	hashFunction := hashFunctionProvider.NewHashFunction()
	log.Infof("Searching rainbow table %s for %d randomly generated passwords...", rainbowTable.Name, numPasswords)

	for i := 0; i < int(numPasswords); i++ {
		password := util.RandomString(&rainbowTable.CharacterSet, rainbowTable.PasswordLength)
		hashedPassword := fmt.Sprintf("%x", hashFunction.Apply(password))
		log.Infof("searching for %s which should reverse to %s", hashedPassword, password)

		_, err := http.Post(
			fmt.Sprintf("%s/api/rainbow-table/%d/search?hash=%s", hashbashHost, rainbowTable.ID, hashedPassword),
			"application/json",
			nil,
		)

		if err != nil {
			log.Errorf("Search request failed: %s", err)
			os.Exit(1)
		}
	}
}

func newSearchSubcommand() *cobra.Command {
	searchCommand := &cobra.Command{
		Use:   "search <rainbow-table-id> <num-passwords>",
		Short: "Generate random passwords and searches a provided rainbow table for their hashes",
		Run:   searchCommandFn,
		Args: cobra.ExactArgs(2),
	}

	searchCommand.Flags().String("hashbash-host", "http://localhost", "Hashbash server to make API requests to")
	viper.BindPFlags(searchCommand.Flags())

	return searchCommand
}
