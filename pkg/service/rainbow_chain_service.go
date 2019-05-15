package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
)

type RainbowChainService interface {
	CreateRainbowChains(int16, []model.RainbowChain) error
	CountChainsForRainbowTable(int16) int64
}

type MySQLRainbowChainService struct {
	databaseClient *gorm.DB
}

func NewRainbowChainService(db *gorm.DB) RainbowChainService {
	return &MySQLRainbowChainService{databaseClient: db}
}

func (service *MySQLRainbowChainService) CreateRainbowChains(rainbowTableId int16, rainbowChains []model.RainbowChain) error {
	sort.Slice(rainbowChains, func(i, j int) bool {
		return rainbowChains[i].EndHash < rainbowChains[j].EndHash
	})

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(fmt.Sprintf("INSERT IGNORE INTO %s (rainbowTableId, startPlaintext, endHash) VALUES ", model.RainbowChain{}.TableName()))
	queryBuilder.WriteString(fmt.Sprintf("(%d, '%s', '%s')", rainbowTableId, rainbowChains[0].StartPlaintext, rainbowChains[0].EndHash))

	for _, r := range rainbowChains[1:] {
		queryBuilder.WriteString(fmt.Sprintf(", (%d, '%s', '%s')", rainbowTableId, r.StartPlaintext, r.EndHash))
	}

	return service.databaseClient.
		Exec(queryBuilder.String()).
		Error
}

func (service *MySQLRainbowChainService) CountChainsForRainbowTable(rainbowTableId int16) int64 {
	var finalChainCount int64
	service.databaseClient.
		Model(&model.RainbowChain{}).
		Where("rainbowTableId = ?", rainbowTableId).
		Count(&finalChainCount)

	return finalChainCount
}
