package userrepo

import (
	"database/sql"
	"github.com/Abdulsametileri/messaging-service/helpers"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"testing"
)

func mockDB() (*gorm.DB, sqlmock.Sqlmock) {
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

func TestCreateUser(t *testing.T) {
	gormDB, mock := mockDB()
	_, _ = gormDB, mock

	t.Run("Create user", func(t *testing.T) {
		user := models.User{
			UserName: "samet123",
			Password: helpers.Sha256String("123456"),
		}

		userRepo := NewUserRepo(gormDB)

		const sqlInsert = `INSERT INTO "user" ("created_at","updated_at","deleted_at","user_name","password","muted_user_ids") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`

		mock.ExpectBegin()
		mock.ExpectQuery(sqlInsert).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.UserName, user.Password, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := userRepo.CreateUser(&user)
		assert.Nil(t, err)
	})
}

func TestGetUserById(t *testing.T) {
	gormDB, mock := mockDB()
	_, _ = gormDB, mock

	t.Run("GetUserById", func(t *testing.T) {
		expected := &models.User{
			UserName: "abdulsamet",
		}

		userRepo := NewUserRepo(gormDB)

		const sql = `SELECT * FROM "user" WHERE "user"."id" = $1 AND "user"."deleted_at" IS NULL ORDER BY "user"."id" LIMIT 1`

		mock.
			ExpectQuery(sql).
			WillReturnRows(
				sqlmock.NewRows([]string{"user_name"}).
					AddRow("abdulsamet"))

		result, err := userRepo.GetUserByID(1)

		assert.EqualValues(t, expected, result)
		assert.Nil(t, err)
	})
}

func TestGetUser(t *testing.T) {
	gormDB, mock := mockDB()
	_, _ = gormDB, mock

	t.Run("GetUser", func(t *testing.T) {
		expected := &models.User{
			UserName: "abdulsamet",
			Password: "123456",
		}

		userRepo := NewUserRepo(gormDB)

		const sql = `SELECT * FROM "user" WHERE (user_name = $1 AND password = $2) AND 
					"user"."deleted_at" IS NULL ORDER BY "user"."id" LIMIT 1`

		mock.
			ExpectQuery(sql).
			WithArgs(expected.UserName, expected.Password).
			WillReturnRows(
				sqlmock.NewRows([]string{"user_name", "password"}).
					AddRow(expected.UserName, expected.Password))

		result, err := userRepo.GetUser(expected.UserName, expected.Password)

		assert.EqualValues(t, expected, result)
		assert.Nil(t, err)
	})
}

func TestGetUserByUserName(t *testing.T) {
	gormDB, mock := mockDB()
	_, _ = gormDB, mock

	t.Run("GetUserByUserName", func(t *testing.T) {
		expected := models.User{
			UserName: "abdulsamet",
			Password: "123456",
		}

		userRepo := NewUserRepo(gormDB)

		const sql = `SELECT * FROM "user" WHERE user_name = $1 AND 
					"user"."deleted_at" IS NULL ORDER BY "user"."id" LIMIT 1`

		mock.
			ExpectQuery(sql).
			WithArgs(expected.UserName).
			WillReturnRows(
				sqlmock.NewRows([]string{"user_name", "password"}).
					AddRow(expected.UserName, expected.Password))

		result, err := userRepo.GetUserByUserName(expected.UserName)

		assert.EqualValues(t, expected, result)
		assert.Nil(t, err)
	})
}

func TestGetUserList(t *testing.T) {
	gormDB, mock := mockDB()
	_, _ = gormDB, mock

	t.Run("Get user list with empty mutated user ids", func(t *testing.T) {
		userId := 1

		rows := sqlmock.
			NewRows([]string{"id", "user_name", "password"}).
			AddRow(1, "abdulsamet", "123456").
			AddRow(2, "anotherperson", "968574")

		userRepo := NewUserRepo(gormDB)

		const sql = `SELECT * FROM "user" WHERE id <> $1 AND "user"."deleted_at" IS NULL`

		mock.ExpectQuery(sql).
			WithArgs(userId).
			WillReturnRows(rows)

		result, err := userRepo.GetUserList(userId, "")

		assert.Nil(t, err)
		assert.Len(t, result, 2)

		assert.Equal(t, result[0].ID, 1)
		assert.Equal(t, result[1].ID, 2)

		assert.Equal(t, result[0].UserName, "abdulsamet")
		assert.Equal(t, result[1].UserName, "anotherperson")

		assert.Equal(t, result[0].Password, "123456")
		assert.Equal(t, result[1].Password, "968574")
	})

}
