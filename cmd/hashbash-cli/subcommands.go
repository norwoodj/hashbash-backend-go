package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/norwoodj/hashbash-backend-go/pkg/api_model"
	"github.com/norwoodj/hashbash-backend-go/pkg/rainbow"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func parseCliArgs(args []string) (int, int) {
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

	return int(rainbowTableId), int(numPasswords)
}

func retrieveRainbowTable(hashbashHost string, rainbowTableId int) api_model.RainbowTable {
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

	var rainbowTable api_model.RainbowTable
	err = json.Unmarshal(body, &rainbowTable)

	if err != nil {
		log.Errorf("Failed to parse json rainbow table response body: %s", err)
		os.Exit(1)
	}

	return rainbowTable
}

func getHashFunction(rainbowTable api_model.RainbowTable) rainbow.HashFunction {
	hashFunctionProvider, err := rainbow.GetHashFunctionProvider(rainbowTable.HashFunction)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	return hashFunctionProvider.NewHashFunction()
}

func submitSearchRequests(
	hashbashHost string,
	rainbowTable api_model.RainbowTable,
	numPasswords int,
) {
	log.Infof("Searching rainbow table %s for %d randomly generated passwords...", rainbowTable.Name, numPasswords)
	hashFunction := getHashFunction(rainbowTable)
	randomStringGenerator := rainbow.NewRandomStringGenerator(1)

	for i := 0; i < int(numPasswords); i++ {
		password := randomStringGenerator.NewRandomString(rainbowTable.CharacterSet, rainbowTable.PasswordLength)
		hashedPassword := hex.EncodeToString(hashFunction.Apply(password))
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

func searchCommandFn(_ *cobra.Command, args []string) {
	hashbashHost := viper.GetString("hashbash-host")
	rainbowTableId, numPasswords := parseCliArgs(args)
	rainbowTable := retrieveRainbowTable(hashbashHost, rainbowTableId)

	submitSearchRequests(hashbashHost, rainbowTable, numPasswords)
}

func newSearchSubcommand() *cobra.Command {
	searchCommand := &cobra.Command{
		Use:   "search <rainbow-table-id> <num-passwords>",
		Short: "Generate random passwords and searches a provided rainbow table for their hashes",
		Run:   searchCommandFn,
		Args:  cobra.ExactArgs(2),
	}

	searchCommand.Flags().String("hashbash-host", "http://localhost", "Hashbash server to make API requests to")
	viper.BindPFlags(searchCommand.Flags())

	return searchCommand
}
