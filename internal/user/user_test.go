package user_test

import (
	"go-post/internal/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserStruct(t *testing.T) {
	user := user.User{
		Id:       1,
		Email:    "test@gmail.com",
		Password: "jdbvj232342",
	}

	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "test@gmail.com", user.Email)
	assert.Equal(t, "jdbvj232342", user.Password)
}
