package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

func MockDB() (*gorm.DB, sqlmock.Sqlmock) {
	var mock sqlmock.Sqlmock
	var db *sql.DB
	var err error

	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		log.Fatalf("can't create sqlmock: %s", err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gdb, err := gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})

	if err != nil {
		log.Fatalf("can't open gorm connection: %s", err)
	}

	return gdb, mock
}
