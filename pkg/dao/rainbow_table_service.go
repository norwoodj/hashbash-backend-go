package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
	"time"
)

type RainbowTableService interface {
	CountRainbowTables() int64
	CreateRainbowTable(*model.RainbowTable) (*model.RainbowTable, error)
	DeleteRainbowTableById(int16) error
	FindRainbowTableById(int16) model.RainbowTable
	FindRainbowTableByName(string) model.RainbowTable
	IncrementRainbowTableChainsGenerated(int16, int64) error
	ListRainbowTables(PageConfig) []model.RainbowTable
	UpdateRainbowTableStatus(int16, string) error
	UpdateRainbowTableStatusAndFinalChainCount(int16, string, int64) error
	UpdateRainbowTableStatusAndGenerateStarted(int16, string) error
}

type DbRainbowTableService struct {
	databaseClient *gorm.DB
}

func NewRainbowTableService(db *gorm.DB) RainbowTableService {
	return &DbRainbowTableService{databaseClient: db}
}

func (service *DbRainbowTableService) CreateRainbowTable(rainbowTable *model.RainbowTable) (*model.RainbowTable, error) {
	if service.FindRainbowTableByName(rainbowTable.Name).Name != "" {
		return nil, RainbowTableExistsError{Name: rainbowTable.Name}
	}

	err := service.databaseClient.
		Save(rainbowTable).
		Error

	return rainbowTable, err
}

func (service *DbRainbowTableService) CountRainbowTables() int64 {
	var rainbowTableCount int64
	service.databaseClient.
		Model(&model.RainbowTable{}).
		Count(&rainbowTableCount)

	return rainbowTableCount
}

func (service *DbRainbowTableService) ListRainbowTables(pageConfig PageConfig) []model.RainbowTable {
	rainbowTables := make([]model.RainbowTable, 0)

	applyPaging(service.databaseClient, pageConfig).
		Find(&rainbowTables)

	return rainbowTables
}

func (service *DbRainbowTableService) FindRainbowTableById(rainbowTableId int16) model.RainbowTable {
	var rainbowTable model.RainbowTable
	service.databaseClient.
		Where("id = ?", rainbowTableId).
		First(&rainbowTable)

	return rainbowTable
}

func (service *DbRainbowTableService) FindRainbowTableByName(name string) model.RainbowTable {
	var rainbowTable model.RainbowTable
	service.databaseClient.
		Where("name = ?", name).
		First(&rainbowTable)

	return rainbowTable
}

func (service *DbRainbowTableService) DeleteRainbowTableById(id int16) error {
	var rainbowTable model.RainbowTable
	service.databaseClient.
		Where("id = ?", id).
		First(&rainbowTable)

	if rainbowTable.Name == "" {
		return RainbowTableNotExistsError{ID: id}
	}

	return service.databaseClient.
		Delete(rainbowTable).
		Error
}

func (service *DbRainbowTableService) UpdateRainbowTableStatus(id int16, status string) error {
	return service.databaseClient.
		Model(&model.RainbowTable{ID: id}).
		Update("status", status).
		Error
}

func (service *DbRainbowTableService) UpdateRainbowTableStatusAndGenerateStarted(id int16, status string) error {
	return service.databaseClient.
		Model(&model.RainbowTable{ID: id}).
		Updates(map[string]interface{}{
			"status":           status,
			"generate_started": time.Now(),
		}).Error
}

func (service *DbRainbowTableService) UpdateRainbowTableStatusAndFinalChainCount(id int16, status string, finalChainCount int64) error {
	return service.databaseClient.
		Model(&model.RainbowTable{ID: id}).
		Updates(map[string]interface{}{
			"final_chain_count":  finalChainCount,
			"generate_completed": time.Now(),
			"status":             status,
		}).Error
}

func (service *DbRainbowTableService) IncrementRainbowTableChainsGenerated(id int16, chainsGenerated int64) error {
	return service.databaseClient.
		Model(&model.RainbowTable{ID: id}).
		Update("chains_generated", gorm.Expr("chains_generated + ?", chainsGenerated)).
		Error
}
