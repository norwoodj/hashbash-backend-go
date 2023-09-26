package dao

import (
	"gorm.io/gorm"
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

	orderClause := pageConfig.SortKey
	if pageConfig.Descending {
		orderClause += " DESC"
	}

	return db.Limit(limit).
		Offset(offset).
		Order(orderClause)
}
