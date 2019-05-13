package service

import (
	"github.com/jinzhu/gorm"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
)

type RainbowTableService interface {
	CountRainbowTables() int64
	CreateRainbowTable(*model.RainbowTable) (*model.RainbowTable, error)
	ListRainbowTables(PageConfig) []model.RainbowTable
	FindRainbowTableById(int16) model.RainbowTable
	FindRainbowTableByName(string) model.RainbowTable
}

type MySQLRainbowTableService struct {
	databaseClient *gorm.DB
}

func NewRainbowTableService(db *gorm.DB) RainbowTableService {
	return MySQLRainbowTableService{databaseClient: db}
}

struct RainbowTableExistsError {

}

func (service MySQLRainbowTableService) CreateRainbowTable(rainbowTable *model.RainbowTable) (*model.RainbowTable, error) {
	if service.FindRainbowTableByName(rainbowTable.Name).Name != "" {
		return
	}
	err := service.databaseClient.
		Save(rainbowTable).
		Error

	return rainbowTable, err
}

func (service MySQLRainbowTableService) CountRainbowTables() int64 {
	var rainbowTableCount int64
	service.databaseClient.
		Model(&model.RainbowTable{}).
		Count(&rainbowTableCount)

	return rainbowTableCount
}

func (service MySQLRainbowTableService) ListRainbowTables(pageConfig PageConfig) []model.RainbowTable {
	rainbowTables := make([]model.RainbowTable, 0)

	applyPaging(service.databaseClient, pageConfig).
		Find(&rainbowTables)

	return rainbowTables
}

func (service MySQLRainbowTableService) FindRainbowTableById(rainbowTableId int16) model.RainbowTable {
	var rainbowTable model.RainbowTable
	service.databaseClient.
		Where("id = ?", rainbowTableId).
		First(&rainbowTable)

	return rainbowTable
}

func (service MySQLRainbowTableService) FindRainbowTableByName(name string) model.RainbowTable {
	var rainbowTable model.RainbowTable
	service.databaseClient.
		Where("name = ?", name).
		First(&rainbowTable)

	return rainbowTable
}
