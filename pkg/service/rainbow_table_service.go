package service

import (
	"github.com/jinzhu/gorm"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
)

type RainbowTableService interface {
	CountRainbowTables() int64
	ListRainbowTables(PageConfig) []model.RainbowTable
	FindRainbowTableById(int16) model.RainbowTable
	CreateRainbowTable(*model.RainbowTable) *model.RainbowTable
}

type MySQLRainbowTableService struct {
	DatabaseClient *gorm.DB
}

func NewRainbowTableService(db *gorm.DB) RainbowTableService {
	return MySQLRainbowTableService{DatabaseClient: db}
}

func (service MySQLRainbowTableService) CreateRainbowTable(rainbowTable *model.RainbowTable) *model.RainbowTable {
	service.DatabaseClient.
		Save(rainbowTable)

	return rainbowTable
}

func (service MySQLRainbowTableService) CountRainbowTables() int64 {
	var rainbowTableCount int64
	service.DatabaseClient.
		Model(&model.RainbowTable{}).
		Count(&rainbowTableCount)

	return rainbowTableCount
}

func (service MySQLRainbowTableService) ListRainbowTables(pageConfig PageConfig) []model.RainbowTable {
	rainbowTables := make([]model.RainbowTable, 0)

	applyPaging(service.DatabaseClient, pageConfig).
		Find(&rainbowTables)

	return rainbowTables
}

func (service MySQLRainbowTableService) FindRainbowTableById(rainbowTableId int16) model.RainbowTable {
	var rainbowTable model.RainbowTable
	service.DatabaseClient.
		Where("id = ?", rainbowTableId).
		First(&rainbowTable)

	return rainbowTable
}
