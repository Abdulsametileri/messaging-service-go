package controllers

import "github.com/Abdulsametileri/messaging-service/models"

type messageSvc struct{}

func (ms *messageSvc) GetMessages(senderUserId, receiverUserId int) ([]models.Message, error) {
	ret := make([]models.Message, 2)
	ret[0] = models.Message{
		BaseModel:  models.BaseModel{ID: 1},
		SenderId:   1,
		ReceiverId: 2,
		Text:       "1-2 ye mesaj attı",
	}
	ret[1] = models.Message{
		BaseModel:  models.BaseModel{ID: 2},
		SenderId:   2,
		ReceiverId: 1,
		Text:       "2-1 e mesaj attı",
	}

	return ret, nil
}

func (ms *messageSvc) CreateMessage(message *models.Message) error {
	return nil
}
