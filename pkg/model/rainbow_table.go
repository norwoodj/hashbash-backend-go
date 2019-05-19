package model

import (
	"time"
)

const StatusQueued = "QUEUED"
const StatusStarted = "STARTED"
const StatusFailed = "FAILED"
const StatusCompleted = "COMPLETED"

const StatusFound = "FOUND"
const StatusNotFound = "NOT_FOUND"

type RainbowTable struct {
	ID                int16      `gorm:"primary_key,column:id"`
	Name              string     `gorm:"column:name"`
	ChainLength       int64      `gorm:"column:chainLength"`
	ChainsGenerated   int64      `gorm:"column:chainsGenerated"`
	CharacterSet      string     `gorm:"column:characterSet"`
	FinalChainCount   int64      `gorm:"column:finalChainCount"`
	HashFunction      string     `gorm:"column:hashFunction"`
	NumChains         int64      `gorm:"column:numChains"`
	PasswordLength    int64      `gorm:"column:passwordLength"`
	Status            string     `gorm:"column:status"`
	GenerateStarted   *time.Time `gorm:"column:generateStarted"`
	GenerateCompleted *time.Time `gorm:"column:generateCompleted"`
	CreatedAt         *time.Time `gorm:"column:created"`
	UpdatedAt         *time.Time `gorm:"column:lastUpdated"`
}

func (RainbowTable) TableName() string {
	return RainbowTableTableName
}
