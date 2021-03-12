package repository

import (
	"database/sql"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var _ = Describe("Auth Repository", func() {
	var repo *AuthRepository
	var mock sqlmock.Sqlmock

	BeforeEach(func() {
		var db *sql.DB
		var err error

		db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		Expect(err).ShouldNot(HaveOccurred())

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{
			NamingStrategy: schema.NamingStrategy{SingularTable: true},
		})

		Expect(err).ShouldNot(HaveOccurred())

		repo = &AuthRepository{Repository{db: gdb}}
	})
	AfterEach(func() {
		err := mock.ExpectationsWereMet() // make sure all expectations were met
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("user", func() {
		It("Exist User", func() {
			user := existMockUser(mock)
			Expect(user.ID).Should(Equal(1))

			existUser, err := repo.UserExist(&user)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(existUser.ID).Should(Equal(user.ID))
			Expect(existUser.UserName).Should(Equal(user.UserName))
			Expect(existUser.Password).Should(Equal(user.Password))
		})
		It("Create User", func() {
			const newId = 1

			user := createMockUser(mock, newId)
			Expect(user.ID).Should(BeZero())

			err := repo.CreateUser(&user)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(user.ID).Should(BeEquivalentTo(newId))
		})
	})
})

func existMockUser(mock sqlmock.Sqlmock) models.User {
	now := time.Now()

	user := models.User{
		BaseModel: models.BaseModel{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserName: "samet123",
		Password: "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92",
	}

	rows := sqlmock.
		NewRows([]string{"id", "created_at", "updated_at", "user_name", "password"}).
		AddRow(user.ID, user.CreatedAt, user.UpdatedAt, user.UserName, user.Password)

	const sqlSelectOne = `SELECT * FROM "user" WHERE (user_name = $1 AND password = $2) 
			AND "user"."deleted_at" IS NULL AND "user"."id" = $3 ORDER BY "user"."id" LIMIT 1`

	mock.ExpectQuery(sqlSelectOne).
		WithArgs(user.UserName, user.Password, sqlmock.AnyArg()).
		WillReturnRows(rows)

	return user
}

func createMockUser(mock sqlmock.Sqlmock, newId int) models.User {
	user := models.User{
		UserName: "samet123",
		Password: "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92",
	}

	const sqlInsert = `INSERT INTO "user" 
					("created_at","updated_at","deleted_at","user_name","password") 
					VALUES ($1,$2,$3,$4,$5) RETURNING "id"`

	mock.ExpectBegin()

	mock.ExpectQuery(sqlInsert).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.UserName, user.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newId))

	mock.ExpectCommit()

	return user
}
