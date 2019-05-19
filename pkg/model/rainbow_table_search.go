package model

import (
	"time"
)

const SearchFound = "FOUND"
const SearchNotFound = "NOT_FOUND"

type RainbowTableSearch struct {
	ID              int64      `gorm:"primary_key,column:id"`
	RainbowTableId  int16      `gorm:"column:rainbowTableId"`
	Hash            string     `gorm:"column:hash"`
	Status          string     `gorm:"column:status"`
	Password        string     `gorm:"password"`
	SearchStarted   *time.Time `gorm:"column:searchStarted"`
	SearchCompleted *time.Time `gorm:"column:searchCompleted"`
	CreatedAt       *time.Time `gorm:"column:created"`
	UpdatedAt       *time.Time `gorm:"column:lastUpdated"`
}

func (RainbowTableSearch) TableName() string {
	return RainbowTableSearchTableName
}
