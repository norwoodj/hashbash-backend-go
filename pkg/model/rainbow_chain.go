package model

type RainbowChain struct {
	RainbowTableId int16  `gorm:"column:rainbowTableId"`
	StartPlaintext string `gorm:"column:startPlaintext"`
	EndHash        string `gorm:"column:endHash"`
}

func (RainbowChain) TableName() string {
	return RainbowChainTableName
}
