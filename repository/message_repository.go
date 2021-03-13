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

func (p *MessageRepository) CreateMessage(message *models.Message) error {
	err := p.db.Model(&models.Message{}).Create(&message).Error
	return err
}
