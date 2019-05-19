package api_model

import (
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
	"time"
)

type RainbowTable struct {
	ID                int16      `json:"id"`
	Name              string     `json:"name"`
	ChainLength       int64      `json:"chainLength"`
	ChainsGenerated   int64      `json:"chainsGenerated"`
	CharacterSet      string     `json:"characterSet"`
	FinalChainCount   int64      `json:"finalChainCount"`
	HashFunction      string     `json:"hashFunction"`
	NumChains         int64      `json:"numChains"`
	PasswordLength    int64      `json:"passwordLength"`
	Status            string     `json:"status"`
	GenerateStarted   *time.Time `json:"generateStarted"`
	GenerateCompleted *time.Time `json:"generateCompleted"`
	GenerationTime    *float64   `json:"generationTime"`
	CreatedAt         *time.Time `json:"created"`
}

func ConvertRainbowTableToApiModel(rainbowTable model.RainbowTable) RainbowTable {
	var generationTimeResult *float64
	if rainbowTable.GenerateCompleted != nil {
		generationTime := rainbowTable.GenerateCompleted.Sub(*rainbowTable.GenerateStarted).Seconds()
		generationTimeResult = &generationTime
	}

	return RainbowTable{
		ID:                rainbowTable.ID,
		Name:              rainbowTable.Name,
		ChainLength:       rainbowTable.ChainLength,
		ChainsGenerated:   rainbowTable.ChainsGenerated,
		CharacterSet:      rainbowTable.CharacterSet,
		FinalChainCount:   rainbowTable.FinalChainCount,
		HashFunction:      rainbowTable.HashFunction,
		NumChains:         rainbowTable.NumChains,
		PasswordLength:    rainbowTable.PasswordLength,
		Status:            rainbowTable.Status,
		GenerateStarted:   rainbowTable.GenerateStarted,
		GenerateCompleted: rainbowTable.GenerateCompleted,
		CreatedAt:         rainbowTable.CreatedAt,
		GenerationTime:    generationTimeResult,
	}
}

func ConvertRainbowTablesToApiModels(rainbowTables []model.RainbowTable) []RainbowTable {
	result := make([]RainbowTable, 0)
	for _, r := range rainbowTables {
		result = append(result, ConvertRainbowTableToApiModel(r))
	}

	return result
}
