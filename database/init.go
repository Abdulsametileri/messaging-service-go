package database

import (
	"fmt"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

var db *gorm.DB

func Setup() *gorm.DB {
	constr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		viper.GetString("DbHost"),
		viper.GetString("DbUser"),
		viper.GetString("DbPass"),
		viper.GetString("DbName"),
		viper.GetString("DbPort"),
	)

	var err error
	db, err = gorm.Open(postgres.Open(constr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}

func Init() {
	db.AutoMigrate(&models.User{})
}
