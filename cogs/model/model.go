package model

type MessageSendReq struct {
	Content string `json:"content"`
	Channel uint64 `json:"channel"`
}
