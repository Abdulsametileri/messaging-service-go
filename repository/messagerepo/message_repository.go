package messagerepo

import (
	"github.com/Abdulsametileri/messaging-service/models"
	"gorm.io/gorm"
)

type Repo interface {
	CreateMessage(message *models.Message) error
	GetMessages(senderUserId, receiverUserId int) ([]models.Message, error)
}

type messageRepo struct {
	db *gorm.DB
}

func NewMessageRepo(db *gorm.DB) Repo {
	return &messageRepo{db: db}
}

func (mr *messageRepo) CreateMessage(message *models.Message) error {
	return mr.db.Model(&models.Message{}).Create(&message).Error
}

func (mr *messageRepo) GetMessages(senderUserId, receiverUserId int) ([]models.Message, error) {
	messages := make([]models.Message, 0)

	err := mr.db.Model(&models.Message{}).
		Joins("Sender").
		Joins("Receiver").
		Where("sender_id = ? AND receiver_id = ?", senderUserId, receiverUserId).
		Or("sender_id = ? AND receiver_id = ?", receiverUserId, senderUserId).
		Order("created_at DESC").
		Find(&messages).
		Error

	return messages, err
}
