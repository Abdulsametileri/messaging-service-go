package api

import (
	"github.com/Abdulsametileri/messaging-service/helpers"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/Abdulsametileri/messaging-service/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var _ = Describe("Auth API Endpoints Test", func() {
	var mock sqlmock.Sqlmock

	BeforeEach(func() {
		var gdb *gorm.DB
		gdb, mock = mockDB()

		repository.Setup(repository.LogRepository{}, gdb)
		repository.Setup(repository.AuthRepository{}, gdb)

		mockGIN()
	})

	AfterEach(func() {
		err := mock.ExpectationsWereMet()
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("register", func() {
		It("added new user successfully", func() {
			const newId = 1
			mockCreateUser(mock, newId)

			payload := gin.H{
				"user_name": "samet123",
				"password":  "123456",
			}

			w, c := performRequest("POST", "api/v1/register", payload)
			c.Set("test", true)

			register(c)

			res := fillResponse(w)

			Expect(res.Message).ShouldNot(Equal(ErrUserAlreadyExist.Error()))
			Expect(res.Code).Should(Equal(http.StatusOK))
		})

		It("user is already exist", func() {
			mockExistUser(mock)

			payload := gin.H{
				"user_name": "samet123",
				"password":  "123456",
			}

			w, c := performRequest("POST", "api/v1/register", payload)

			register(c)

			res := fillResponse(w)

			Expect(res.Message).Should(Equal(ErrUserAlreadyExist.Error()))
			Expect(res.Code).Should(Equal(http.StatusBadRequest))
		})
	})

	Context("login", func() {
		It("user does not exist", func() {
			payload := gin.H{
				"user_name": "samet123",
				"password":  "123456",
			}

			w, c := performRequest("POST", "api/v1/login", payload)
			c.Set("test", true)
			login(c)

			res := fillResponse(w)

			Expect(res.Code).Should(Equal(http.StatusBadRequest))
			Expect(res.Message).Should(Equal(ErrUserDoesNotExist.Error()))
		})
		It("successfully login", func() {
			mockGetUser(mock)

			payload := gin.H{
				"user_name": "samet123",
				"password":  "123456",
			}

			w, c := performRequest("POST", "api/v1/login", payload)

			login(c)

			res := fillResponse(w)

			Expect(res.Message).Should(Equal(""))
			Expect(res.Code).Should(Equal(http.StatusOK))
		})
	})
})

func mockExistUser(mock sqlmock.Sqlmock) {
	now := time.Now()

	user := models.User{
		BaseModel: models.BaseModel{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserName: "samet123",
		Password: helpers.Sha256String("123456"),
	}

	rows := sqlmock.
		NewRows([]string{"id", "created_at", "updated_at", "user_name", "password"}).
		AddRow(user.ID, user.CreatedAt, user.UpdatedAt, user.UserName, user.Password)

	const sqlSelectOne = `SELECT * FROM "user" WHERE user_name = $1
			AND "user"."deleted_at" IS NULL ORDER BY "user"."id" LIMIT 1`

	mock.ExpectQuery(sqlSelectOne).
		WithArgs(user.UserName).
		WillReturnRows(rows)
}

func mockGetUser(mock sqlmock.Sqlmock) models.User {
	now := time.Now()

	user := models.User{
		BaseModel: models.BaseModel{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserName: "samet123",
		Password: helpers.Sha256String("123456"),
	}

	rows := sqlmock.
		NewRows([]string{"id", "created_at", "updated_at", "user_name", "password"}).
		AddRow(user.ID, user.CreatedAt, user.UpdatedAt, user.UserName, user.Password)

	const sqlSelectOne = `SELECT * FROM "user" WHERE (user_name = $1 AND password = $2)
			AND "user"."deleted_at" IS NULL ORDER BY "user"."id" LIMIT 1`

	mock.ExpectQuery(sqlSelectOne).
		WithArgs(user.UserName, user.Password).
		WillReturnRows(rows)

	return user
}

func mockCreateUser(mock sqlmock.Sqlmock, newId int) models.User {
	user := models.User{
		UserName: "samet123",
		Password: helpers.Sha256String("123456"),
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
