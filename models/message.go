package models

type Message struct {
	BaseModel
	Sender     User   `json:"sender"`
	SenderId   int    `json:"sender_id"`
	Receiver   User   `json:"receiver"`
	ReceiverId int    `json:"receiver_id"`
	Text       string `json:"text"`
}
