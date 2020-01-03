package dao

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
	FindChainByTableIdAndEndHashIn(int16, []string) []model.RainbowChain
}

type DbRainbowChainService struct {
	databaseClient *gorm.DB
}

func NewRainbowChainService(db *gorm.DB) RainbowChainService {
	return &DbRainbowChainService{databaseClient: db}
}

func (service *DbRainbowChainService) CreateRainbowChains(rainbowTableId int16, rainbowChains []model.RainbowChain) error {
	sort.Slice(rainbowChains, func(i, j int) bool {
		return rainbowChains[i].EndHash < rainbowChains[j].EndHash
	})

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(fmt.Sprintf("INSERT INTO %s (rainbow_table_id, start_plaintext, end_hash) VALUES ", model.RainbowChain{}.TableName()))
	queryBuilder.WriteString(fmt.Sprintf("(%d, '%s', '%s')", rainbowTableId, rainbowChains[0].StartPlaintext, rainbowChains[0].EndHash))

	for _, r := range rainbowChains[1:] {
		queryBuilder.WriteString(fmt.Sprintf(", (%d, '%s', '%s')", rainbowTableId, r.StartPlaintext, r.EndHash))
	}

	queryBuilder.WriteString(" ON CONFLICT(rainbow_table_id, end_hash) DO NOTHING")
	return service.databaseClient.
		Exec(queryBuilder.String()).
		Error
}

func (service *DbRainbowChainService) CountChainsForRainbowTable(rainbowTableId int16) int64 {
	var finalChainCount int64
	service.databaseClient.
		Model(&model.RainbowChain{}).
		Where("rainbow_table_id = ?", rainbowTableId).
		Count(&finalChainCount)

	return finalChainCount
}

func (service *DbRainbowChainService) FindChainByTableIdAndEndHashIn(rainbowTableId int16, endHashes []string) []model.RainbowChain {
	var rainbowChains []model.RainbowChain
	service.databaseClient.
		Model(&model.RainbowChain{}).
		Where("rainbow_table_id = ? AND end_hash IN (?)", rainbowTableId, endHashes).
		Scan(&rainbowChains)

	return rainbowChains
}
