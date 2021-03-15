package messageservice

import (
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/Abdulsametileri/messaging-service/repository/messagerepo"
)

type MessageService interface {
	GetMessages(senderUserId, receiverUserId int) ([]models.Message, error)
	CreateMessage(message *models.Message) error
}

type messageService struct {
	Repo messagerepo.Repo
}

func NewMessageService(repo messagerepo.Repo) MessageService {
	return &messageService{Repo: repo}
}

func (ms *messageService) CreateMessage(message *models.Message) error {
	return ms.Repo.CreateMessage(message)
}

func (ms *messageService) GetMessages(senderUserId, receiverUserId int) ([]models.Message, error) {
	return ms.Repo.GetMessages(senderUserId, receiverUserId)
}
