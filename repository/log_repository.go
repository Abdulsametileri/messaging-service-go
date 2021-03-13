package repository

import (
	"github.com/Abdulsametileri/messaging-service/models"
	"gorm.io/gorm"
)

var logRepo *LogRepository

type LogRepository struct {
	Repository
}

func (a LogRepository) Setup(db *gorm.DB) {
	logRepo = &LogRepository{Repository{db: db}}
}

func GetLogRepository() *LogRepository {
	return logRepo
}

func (a *LogRepository) CreateLog(log *models.Log) error {
	return logRepo.db.Model(&models.Log{}).Create(&log).Error
}
