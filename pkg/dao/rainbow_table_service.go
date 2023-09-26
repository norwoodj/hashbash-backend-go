package dao

import (
	"errors"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
	"gorm.io/gorm"
	"time"
)

type RainbowTableService interface {
	CountRainbowTables() (int64, error)
	CreateRainbowTable(*model.RainbowTable) (*model.RainbowTable, error)
	DeleteRainbowTableById(int16) error
	FindRainbowTableById(int16) (model.RainbowTable, error)
	FindRainbowTableByName(string) (model.RainbowTable, error)
	IncrementRainbowTableChainsGenerated(int16, int64) error
	ListRainbowTables(PageConfig) ([]model.RainbowTable, error)
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
	_, err := service.FindRainbowTableByName(rainbowTable.Name)
	if err == nil {
		return nil, RainbowTableExistsError{Name: rainbowTable.Name}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	err = service.databaseClient.
		Save(rainbowTable).
		Error

	return rainbowTable, err
}

func (service *DbRainbowTableService) CountRainbowTables() (int64, error) {
	var rainbowTableCount int64
	err := service.databaseClient.
		Model(&model.RainbowTable{}).
		Count(&rainbowTableCount).
		Error

	return rainbowTableCount, err
}

func (service *DbRainbowTableService) ListRainbowTables(pageConfig PageConfig) ([]model.RainbowTable, error) {
	rainbowTables := make([]model.RainbowTable, 0)

	err := applyPaging(service.databaseClient, pageConfig).
		Find(&rainbowTables).
		Error

	return rainbowTables, err
}

func (service *DbRainbowTableService) FindRainbowTableById(rainbowTableId int16) (model.RainbowTable, error) {
	var rainbowTable model.RainbowTable
	err := service.databaseClient.
		Where("id = ?", rainbowTableId).
		First(&rainbowTable).
		Error

	return rainbowTable, err
}

func (service *DbRainbowTableService) FindRainbowTableByName(name string) (model.RainbowTable, error) {
	var rainbowTable model.RainbowTable
	err := service.databaseClient.
		Where("name = ?", name).
		First(&rainbowTable).
		Error

	return rainbowTable, err
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
