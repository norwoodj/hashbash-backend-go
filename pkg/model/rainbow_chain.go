package model

type RainbowChain struct {
	RainbowTableId int16  `gorm:"column:rainbowTableId"json:"rainbowTableId"`
	StartPlaintext string `gorm:"column:startPlaintext"json:"startPlaintext"`
	EndHash        string `gorm:"column:endHash"json:"endHash"`
}

func (RainbowChain) TableName() string {
	return RainbowChainTableName
}
