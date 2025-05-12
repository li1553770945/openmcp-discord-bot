package model

type MessageSendReq struct {
	Content string `json:"content"`
	Channel string `json:"channel"`
}
