package logrepo

import (
	"github.com/Abdulsametileri/messaging-service/models"
	"gorm.io/gorm"
)

type Repo interface {
	CreateLog(log *models.Log) error
}

type logRepo struct {
	db *gorm.DB
}

func NewLogRepo(db *gorm.DB) Repo {
	return &logRepo{db: db}
}

func (l *logRepo) CreateLog(log *models.Log) error {
	return l.db.Model(&models.Log{}).Create(&log).Error
}
