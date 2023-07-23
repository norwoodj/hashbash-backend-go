package model

type RainbowChain struct {
	RainbowTableId int16  `gorm:"column:rainbow_table_id"`
	StartPlaintext string `gorm:"column:start_plaintext"`
	EndHash        []byte `gorm:"column:end_hash"`
}

func (RainbowChain) TableName() string {
	return RainbowChainTableName
}
