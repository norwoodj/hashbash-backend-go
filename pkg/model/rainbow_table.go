package model

import (
	"time"
)

type RainbowTable struct {
	ID              int64     `gorm:"primary_key,column:id"json:"id"`
	Name            string    `gorm:"column:name"json:"name"`
	NumChains       int64     `gorm:"column:numChains"json:"numChains"`
	ChainLength     int64     `gorm:"column:chainLength"json:"chainLength"`
	PasswordLength  int64     `gorm:"column:passwordLength"json:"passwordLength"`
	CharacterSet    string    `gorm:"column:characterSet"json:"characterSet"`
	HashFunction    string    `gorm:"column:hashFunction"json:"hashFunction"`
	FinalChainCount int64     `gorm:"column:finalChainCount"json:"finalChainCount"`
	ChainsGenerated int64     `gorm:"column:chainsGenerated"json:"chainsGenerated"`
	Status          string    `gorm:"column:status"json:"status"`
	Created         time.Time `gorm:"column:created"json:"created"`
	LastUpdated     time.Time `gorm:"column:lastUpdated"json:"-"`
}

func (RainbowTable) TableName() string {
	return "rainbow_table"
}
