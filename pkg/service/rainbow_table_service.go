package service

import (
	"github.com/jinzhu/gorm"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
)

type RainbowTableService interface {
	CountRainbowTables() int64
	CreateRainbowTable(*model.RainbowTable) (*model.RainbowTable, error)
	DeleteRainbowTableById(int16) error
	FindRainbowTableById(int16) model.RainbowTable
	FindRainbowTableByName(string) model.RainbowTable
	IncrementRainbowTableChainsGenerated(int16, int64) error
	ListRainbowTables(PageConfig) []model.RainbowTable
	UpdateRainbowTableFinalChainCount(int16, int64) error
	UpdateRainbowTableStatus(int16, string) error
}

type MySQLRainbowTableService struct {
	databaseClient *gorm.DB
}

func NewRainbowTableService(db *gorm.DB) RainbowTableService {
	return &MySQLRainbowTableService{databaseClient: db}
}

func (service *MySQLRainbowTableService) CreateRainbowTable(rainbowTable *model.RainbowTable) (*model.RainbowTable, error) {
	if service.FindRainbowTableByName(rainbowTable.Name).Name != "" {
		return nil, RainbowTableExistsError{Name: rainbowTable.Name}
	}

	err := service.databaseClient.
		Save(rainbowTable).
		Error

	return rainbowTable, err
}

func (service *MySQLRainbowTableService) CountRainbowTables() int64 {
	var rainbowTableCount int64
	service.databaseClient.
		Model(&model.RainbowTable{}).
		Count(&rainbowTableCount)

	return rainbowTableCount
}

func (service *MySQLRainbowTableService) ListRainbowTables(pageConfig PageConfig) []model.RainbowTable {
	rainbowTables := make([]model.RainbowTable, 0)

	applyPaging(service.databaseClient, pageConfig).
		Find(&rainbowTables)

	return rainbowTables
}

func (service *MySQLRainbowTableService) FindRainbowTableById(rainbowTableId int16) model.RainbowTable {
	var rainbowTable model.RainbowTable
	service.databaseClient.
		Where("id = ?", rainbowTableId).
		First(&rainbowTable)

	return rainbowTable
}

func (service *MySQLRainbowTableService) FindRainbowTableByName(name string) model.RainbowTable {
	var rainbowTable model.RainbowTable
	service.databaseClient.
		Where("name = ?", name).
		First(&rainbowTable)

	return rainbowTable
}

func (service *MySQLRainbowTableService) DeleteRainbowTableById(id int16) error {
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

func (service *MySQLRainbowTableService) UpdateRainbowTableFinalChainCount(id int16, finalChainCount int64) error {
	return service.databaseClient.
		Model(&model.RainbowTable{ID: id}).
		Update("finalChainCount", finalChainCount).
		Error
}

func (service *MySQLRainbowTableService) UpdateRainbowTableStatus(id int16, status string) error {
	return service.databaseClient.
		Model(&model.RainbowTable{ID: id}).
		Update("status", status).
		Error
}

func (service *MySQLRainbowTableService) IncrementRainbowTableChainsGenerated(id int16, chainsGenerated int64) error {
	return service.databaseClient.
		Model(&model.RainbowTable{ID: id}).
		Update("chainsGenerated", gorm.Expr("chainsGenerated + ?", chainsGenerated)).
		Error
}
