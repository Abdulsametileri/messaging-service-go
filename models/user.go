package models

type User struct {
	BaseModel
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
