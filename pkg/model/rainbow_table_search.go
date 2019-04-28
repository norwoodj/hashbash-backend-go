package model

import (
	"time"
)

const SearchFound = "FOUND"
const SearchNotFound = "NOT_FOUND"

type RainbowTableSearch struct {
	ID              int64     `gorm:"primary_key,column:id"json:"id"`
	RainbowTableId  int16     `gorm:"column:rainbowTableId"json:"rainbowTableId"`
	Hash            string    `gorm:"column:hash"json:"hash"`
	Status          string    `gorm:"column:status"json:"status"`
	Password        string    `gorm:"password"json:"password"`
	SearchStarted   time.Time `gorm:"column:searchStarted"json:"searchStarted"`
	SearchCompleted time.Time `gorm:"column:searchCompleted"json:"searchCompleted"`
	Created         time.Time `gorm:"column:created"json:"created"`
	LastUpdated     time.Time `gorm:"column:lastUpdated"json:"-"`
}

func (RainbowTableSearch) TableName() string {
	return RainbowTableSearchTableName
}
