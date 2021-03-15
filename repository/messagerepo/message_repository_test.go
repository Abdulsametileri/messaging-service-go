package messagerepo

import (
	"database/sql"
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

func TestGetMessages(t *testing.T) {
	gormDB, mock := mockDB()
	_, _ = gormDB, mock

	t.Run("Get messages", func(t *testing.T) {
		senderId := 1
		receiverId := 2

		rows := sqlmock.
			NewRows([]string{"sender_id", "receiver_id"}).
			AddRow(1, 2).
			AddRow(2, 1).
			AddRow(1, 2).
			AddRow(2, 1)

		messageRepo := NewMessageRepo(gormDB)

		const sql = `SELECT "message"."id","message"."created_at","message"."updated_at","message"."deleted_at","message"."sender_id","message"."receiver_id","message"."text","Sender"."id" AS "Sender__id","Sender"."created_at" AS "Sender__created_at","Sender"."updated_at" AS "Sender__updated_at","Sender"."deleted_at" AS "Sender__deleted_at","Sender"."user_name" AS "Sender__user_name","Sender"."password" AS "Sender__password","Sender"."muted_user_ids" AS "Sender__muted_user_ids","Receiver"."id" AS "Receiver__id","Receiver"."created_at" AS "Receiver__created_at","Receiver"."updated_at" AS "Receiver__updated_at","Receiver"."deleted_at" AS "Receiver__deleted_at","Receiver"."user_name" AS "Receiver__user_name","Receiver"."password" AS "Receiver__password","Receiver"."muted_user_ids" AS "Receiver__muted_user_ids" FROM "message" LEFT JOIN "user" "Sender" ON "message"."sender_id" = "Sender"."id" LEFT JOIN "user" "Receiver" ON "message"."receiver_id" = "Receiver"."id" WHERE ((sender_id = $1 AND receiver_id = $2) OR (sender_id = $3 AND receiver_id = $4)) AND "message"."deleted_at" IS NULL ORDER BY created_at DESC`

		mock.ExpectQuery(sql).
			WithArgs(senderId, receiverId, receiverId, senderId).
			WillReturnRows(rows)

		messages, err := messageRepo.GetMessages(senderId, receiverId)
		assert.Nil(t, err)

		assert.Len(t, messages, 4)
		assert.Equal(t, messages[0].SenderId, senderId)
		assert.Equal(t, messages[0].ReceiverId, receiverId)

		assert.Equal(t, messages[1].SenderId, receiverId)
		assert.Equal(t, messages[1].ReceiverId, senderId)

		assert.Equal(t, messages[2].SenderId, senderId)
		assert.Equal(t, messages[2].ReceiverId, receiverId)

		assert.Equal(t, messages[3].SenderId, receiverId)
		assert.Equal(t, messages[3].ReceiverId, senderId)
	})
}
