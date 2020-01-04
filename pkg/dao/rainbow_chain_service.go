package dao

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
	log "github.com/sirupsen/logrus"
)

type insertIgnoreConflictClauseProvider interface {
	getAfterInsertModifier() string
	getEndingModifier() string
}

type mysqlInsertIgnoreConflictClauseProvider struct{}
type postgresqlInsertIgnoreConflictClauseProvider struct{}

func (mysqlInsertIgnoreConflictClauseProvider) getAfterInsertModifier() string {
	return "IGNORE"
}

func (mysqlInsertIgnoreConflictClauseProvider) getEndingModifier() string {
	return ""
}

func (postgresqlInsertIgnoreConflictClauseProvider) getAfterInsertModifier() string {
	return ""
}

func (postgresqlInsertIgnoreConflictClauseProvider) getEndingModifier() string {
	return "ON CONFLICT(rainbow_table_id, end_hash) DO NOTHING"
}

func getInsertIgnoreConflictClauseProviderForEngine(engine string) (insertIgnoreConflictClauseProvider, error) {
	switch engine {
	case "mysql":
		return mysqlInsertIgnoreConflictClauseProvider{}, nil
	case "postgres":
		return postgresqlInsertIgnoreConflictClauseProvider{}, nil
	case "postgresql":
		return postgresqlInsertIgnoreConflictClauseProvider{}, nil
	default:
		return nil, fmt.Errorf("no engine %s found", engine)
	}
}

type RainbowChainService interface {
	CreateRainbowChains(int16, []model.RainbowChain) error
	CountChainsForRainbowTable(int16) int64
	FindChainByTableIdAndEndHashIn(int16, []string) []model.RainbowChain
}

type DbRainbowChainService struct {
	databaseClient *gorm.DB
	insertIgnoreConflictClauseProvider
}

func NewRainbowChainService(db *gorm.DB, databaseEngine string) RainbowChainService {
	insertIgnoreConflictClauseProvider, err := getInsertIgnoreConflictClauseProviderForEngine(databaseEngine)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	return &DbRainbowChainService{databaseClient: db, insertIgnoreConflictClauseProvider: insertIgnoreConflictClauseProvider}
}

func (service *DbRainbowChainService) CreateRainbowChains(rainbowTableId int16, rainbowChains []model.RainbowChain) error {
	sort.Slice(rainbowChains, func(i, j int) bool {
		return rainbowChains[i].EndHash < rainbowChains[j].EndHash
	})

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(fmt.Sprintf(
		"INSERT %s INTO %s (rainbow_table_id, start_plaintext, end_hash) VALUES ",
		service.getAfterInsertModifier(),
		model.RainbowChain{}.TableName(),
	))

	queryBuilder.WriteString(fmt.Sprintf("(%d, '%s', '%s')", rainbowTableId, rainbowChains[0].StartPlaintext, rainbowChains[0].EndHash))

	for _, r := range rainbowChains[1:] {
		queryBuilder.WriteString(fmt.Sprintf(", (%d, '%s', '%s')", rainbowTableId, r.StartPlaintext, r.EndHash))
	}

	queryBuilder.WriteString(service.getEndingModifier())
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
