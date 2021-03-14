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
		viper.GetString("DB_HOST"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_PORT"),
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

func Migrate() {
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Log{})
	db.AutoMigrate(&models.Message{})
}
