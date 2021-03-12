package repository

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

type IRepository interface {
	Setup(db *gorm.DB)
}

func Setup(repository IRepository, db *gorm.DB) {
	repository.Setup(db)
}
