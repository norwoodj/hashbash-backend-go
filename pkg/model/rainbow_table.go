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
	ChainLength       int64      `gorm:"column:chain_length"`
	ChainsGenerated   int64      `gorm:"column:chains_generated"`
	CharacterSet      string     `gorm:"column:character_set"`
	FinalChainCount   int64      `gorm:"column:final_chain_count"`
	HashFunction      string     `gorm:"column:hash_function"`
	NumChains         int64      `gorm:"column:num_chains"`
	PasswordLength    int64      `gorm:"column:password_length"`
	Status            string     `gorm:"column:status"`
	GenerateStarted   *time.Time `gorm:"column:generate_started"`
	GenerateCompleted *time.Time `gorm:"column:generate_completed"`
	CreatedAt         *time.Time `gorm:"column:created"`
}

func (RainbowTable) TableName() string {
	return RainbowTableTableName
}
