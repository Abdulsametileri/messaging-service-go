package logservice

import (
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/Abdulsametileri/messaging-service/repository/logrepo"
)

type LogService interface {
	CreateLog(log *models.Log) error
}

type logService struct {
	Repo logrepo.Repo
}

func NewLogService(repo logrepo.Repo) LogService {
	return &logService{Repo: repo}
}

func (ls *logService) CreateLog(log *models.Log) error {
	return ls.Repo.CreateLog(log)
}
