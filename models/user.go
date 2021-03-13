package models

import "github.com/lib/pq"

type User struct {
	BaseModel
	UserName     string        `json:"user_name"`
	Password     string        `json:"password"`
	MutedUserIDs pq.Int32Array `json:"muted_user_ids" gorm:"type:integer[]"`
}
