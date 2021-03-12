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
)

var _ = Describe("Repository", func() {
	var repo *Repository
	var mock sqlmock.Sqlmock
	_, _ = repo, mock

	BeforeEach(func() {
		var db *sql.DB
		var err error

		db, mock, err = sqlmock.New(
			sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		Expect(err).ShouldNot(HaveOccurred())
		_ = mock
		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "sqlite",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{
			NamingStrategy: schema.NamingStrategy{SingularTable: true},
		})

		Expect(err).ShouldNot(HaveOccurred())

		repo = &Repository{db: gdb}
	})

	Context("user", func() {
		It("Is User Exist", func() {

		})
		It("Create User", func() {
			user := models.User{
				UserName: "samet123",
				Password: "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92",
			}

			const sqlInsert = `INSERT INTO "user" 
					("created_at","updated_at","deleted_at","user_name","password") 
					VALUES ($1,$2,$3,$4,$5) RETURNING "id"`
			const newId = 2

			mock.ExpectBegin()

			mock.ExpectQuery(sqlInsert).
				WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.UserName, user.Password).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newId))

			mock.ExpectCommit()

			Expect(user.ID).Should(BeZero())

			err := repo.CreateUser(&user)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(user.ID).Should(BeEquivalentTo(newId))
		})
	})
})
