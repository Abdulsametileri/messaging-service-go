package models

type Message struct {
	BaseModel
	SenderId   int    `json:"sender_id"`
	ReceiverId int    `json:"receiver_id"`
	Text       string `json:"text"`
}
