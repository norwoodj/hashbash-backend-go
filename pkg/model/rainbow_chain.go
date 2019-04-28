package model

type RainbowChain struct {
	RainbowTableId int16  `gorm:"primary_key,column:rainbowTableId"json:"rainbowTableId"`
	StartPlaintext string `gorm:"startPlaintext"json:"startPlaintext"`
	EndHash        string `gorm:"primary_key,column:endHash"json:"endHash"`
}

func (RainbowChain) TableName() string {
	return RainbowChainTableName
}
