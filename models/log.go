package models

import "gorm.io/datatypes"

type LogType string

const (
	LogInfo  LogType = "Info"
	LogError         = "Error"
)

type Log struct {
	BaseModel
	Type     LogType        `json:"type"`
	ApiPath  string         `json:"api_path"`
	Request  datatypes.JSON `json:"request_body"`
	Response datatypes.JSON `json:"response_body"`
}
