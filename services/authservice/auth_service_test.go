package authservice

import (
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestCreateJwtToken(t *testing.T) {
	t.Run("Create jwt token", func(t *testing.T) {
		as := NewAuthService()

		user := &models.User{
			BaseModel: models.BaseModel{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: gorm.DeletedAt{},
			},
			UserName: "abdulsamet",
		}

		token, err := as.CreateJwtToken(user)

		assert.Nil(t, err)
		assert.IsType(t, "string", token)
		assert.Len(t, token, 148)
	})
}
