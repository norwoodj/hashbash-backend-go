package model

import (
	"time"
)

const StatusQueued = "QUEUED"
const StatusStarted = "STARTED"

type RainbowTable struct {
	ID              int16     `gorm:"primary_key,column:id"json:"id"`
	Name            string    `gorm:"column:name"json:"name"`
	ChainLength     int64     `gorm:"column:chainLength"json:"chainLength"`
	ChainsGenerated int64     `gorm:"column:chainsGenerated"json:"chainsGenerated"`
	CharacterSet    string    `gorm:"column:characterSet"json:"characterSet"`
	FinalChainCount int64     `gorm:"column:finalChainCount"json:"finalChainCount"`
	HashFunction    string    `gorm:"column:hashFunction"json:"hashFunction"`
	NumChains       int64     `gorm:"column:numChains"json:"numChains"`
	PasswordLength  int64     `gorm:"column:passwordLength"json:"passwordLength"`
	Status          string    `gorm:"column:status"json:"status"`
	CreatedAt       time.Time `gorm:"column:created"json:"created"`
	UpdatedAt       time.Time `gorm:"column:lastUpdated"json:"-"`
}

func (RainbowTable) TableName() string {
	return RainbowTableTableName
}
