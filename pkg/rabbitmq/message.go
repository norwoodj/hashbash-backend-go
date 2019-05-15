package rabbitmq

type RainbowTableIdMessage struct {
	RainbowTableId int16 `json:"rainbowTableId"`
}

type RainbowTableSearchRequestMessage struct {
	Hash           string `json:"hash"`
	RainbowTableId int16  `json:"rainbowTableId"`
	SearchId       int64  `json:"searchId"`
}
