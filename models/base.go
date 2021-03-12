package models

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int            `gorm:"primaryKey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `sql:"index" json:"deleted_at" form:"deleted_at"`
}
