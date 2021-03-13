package repository

import (
	"github.com/Abdulsametileri/messaging-service/models"
	"gorm.io/gorm"
)

var messageRepo *MessageRepository

type MessageRepository struct {
	Repository
}

func (a MessageRepository) Setup(db *gorm.DB) {
	messageRepo = &MessageRepository{Repository{db: db}}
}

func GetMessageRepository() *MessageRepository {
	return messageRepo
}

func (p *MessageRepository) GetMessages(senderUserId, receiverUserId int) ([]models.Message, error) {
	messages := make([]models.Message, 0)

	err := p.db.Model(&models.Message{}).
		Joins("Sender").
		Joins("Receiver").
		Where("sender_id = ? AND receiver_id = ?", senderUserId, receiverUserId).
		Or("sender_id = ? AND receiver_id = ?", receiverUserId, senderUserId).
		Order("created_at DESC").
		Find(&messages).
		Error

	return messages, err
}

func (p *MessageRepository) CreateMessage(message *models.Message) error {
	err := p.db.Model(&models.Message{}).Create(&message).Error
	return err
}
