package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/norwoodj/hashbash-backend-go/pkg/database"
)

type PageConfig struct {
	Descending bool
	PageNumber int
	PageSize   int
	SortKey    string
}

func applyPaging(db *gorm.DB, pageConfig PageConfig) *gorm.DB {
	limit := pageConfig.PageSize
	offset := pageConfig.PageNumber * pageConfig.PageSize
	return database.ApplyPaging(db, limit, offset, pageConfig.SortKey, pageConfig.Descending)
}
