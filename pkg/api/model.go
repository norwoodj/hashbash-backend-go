package api

type rainbowTableSearchRequest struct {
	Hash string `json:"hash"`
}

type rainbowTableSearchResponse struct {
	Hash     string `json:"hash"`
	SearchId int64  `json:"searchId"`
	Status   string `json:"status"`
}
