package repository

import (
	"github.com/Abdulsametileri/messaging-service/helpers"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("Auth Repository Test", func() {
	var mock sqlmock.Sqlmock
	var authRepo *AuthRepository

	BeforeEach(func() {
		var gdb *gorm.DB
		gdb, mock = mockDB()

		authRepo = &AuthRepository{Repository{db: gdb}}
	})

	AfterEach(func() {
		err := mock.ExpectationsWereMet()
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("register", func() {
		It("create user success?", func() {
			const newId = 3
			user := models.User{
				UserName: "samet123",
				Password: helpers.Sha256String("123456"),
			}

			const sqlInsert = `INSERT INTO "user" ("created_at","updated_at","deleted_at","user_name","password","muted_user_ids") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`

			mock.ExpectBegin()
			mock.ExpectQuery(sqlInsert).
				WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.UserName, user.Password, sqlmock.AnyArg()).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newId))
			mock.ExpectCommit()

			Expect(user.ID).Should(BeZero())

			err := authRepo.CreateUser(&user)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(user.ID).Should(BeEquivalentTo(newId))
		})
		It("exist user success?", func() {
			expected := models.User{UserName: "samet123"}

			rows := sqlmock.
				NewRows([]string{"user_name"}).
				AddRow(expected.UserName)

			const sqlSelectOne = `SELECT * FROM "user" WHERE user_name = $1
			AND "user"."deleted_at" IS NULL ORDER BY "user"."id" LIMIT 1`

			mock.ExpectQuery(sqlSelectOne).
				WithArgs(expected.UserName).
				WillReturnRows(rows)

			found, err := authRepo.ExistUser(expected.UserName)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(found).Should(BeTrue())
		})
		It("get user success?", func() {
			expected := models.User{
				UserName: "samet123",
				Password: helpers.Sha256String("123456"),
			}

			sqlStr := `SELECT * FROM "user" WHERE (user_name = $1 AND password = $2) AND 
					"user"."deleted_at" IS NULL ORDER BY "user"."id" LIMIT 1`

			mock.
				ExpectQuery(sqlStr).
				WithArgs(expected.UserName, expected.Password).
				WillReturnRows(
					sqlmock.NewRows([]string{"user_name", "password"}).
						AddRow(expected.UserName, expected.Password))

			result, err := authRepo.GetUser(expected.UserName, expected.Password)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(result.UserName).Should(Equal(expected.UserName))
			Expect(result.Password).Should(Equal(expected.Password))
		})
	})
})
