package model

import (
	"time"
)

const SearchFound = "FOUND"
const SearchNotFound = "NOT_FOUND"

type RainbowTableSearch struct {
	ID              int64      `gorm:"primary_key,column:id"`
	RainbowTableId  int16      `gorm:"column:rainbow_table_id"`
	Hash            []byte     `gorm:"column:hash"`
	Status          string     `gorm:"column:status"`
	Password        string     `gorm:"column:password"`
	SearchStarted   *time.Time `gorm:"column:search_started"`
	SearchCompleted *time.Time `gorm:"column:search_completed"`
	CreatedAt       *time.Time `gorm:"column:created"`
}

func (RainbowTableSearch) TableName() string {
	return RainbowTableSearchTableName
}
