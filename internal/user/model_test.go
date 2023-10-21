package user_test

import (
	"go-post/internal/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModelUserStuct(t *testing.T) {
	u := user.User{
		Id:       1,
		Email:    "test@gmail.com",
		Password: "12345",
	}

	assert.Equal(t, 1, u.Id)
	assert.Equal(t, "test@gmail.com", u.Email)
	assert.Equal(t, "12345", u.Password)
}

func TestModelInputUserSignupStuct(t *testing.T) {
	i := user.InputUserSignUp{
		Email:    "test@gmail.com",
		Password: "12345",
	}

	assert.Equal(t, "test@gmail.com", i.Email)
	assert.Equal(t, "12345", i.Password)
}

func TestInputUserLoginStruct(t *testing.T) {
	i := user.InputUserLogin{
		Email:    "test@gmail.com",
		Password: "12345",
	}

	assert.Equal(t, "test@gmail.com", i.Email)
	assert.Equal(t, "12345", i.Password)
}
