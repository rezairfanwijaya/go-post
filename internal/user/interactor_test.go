package user_test

import (
	"errors"
	"go-post/internal/user"
	"go-post/internal/user/mocks"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		Name     string
		User     user.User
		HttpCode int
		Error    error
	}{
		{
			Name: "success",
			User: user.User{
				Email:    "test@gmail.com",
				Password: "djbgj121-232j",
			},
			HttpCode: http.StatusOK,
		},
		{
			Name:     "failed",
			User:     user.User{},
			HttpCode: http.StatusInternalServerError,
			Error:    errors.New("failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)
			interactor := user.NewInteractor(repo)

			repo.On("Save", tc.User).Return(tc.User, tc.Error)
			res, httpCode, err := interactor.CreateUser(tc.User)
			assert.Equal(t, tc.Error, err)
			assert.Equal(t, tc.User, res)
			assert.Equal(t, tc.HttpCode, httpCode)
		})
	}
}

func TestValidateUser(t *testing.T) {
	testCases := []struct {
		Name    string
		UserId  int
		User    user.User
		IsValid bool
		Error   error
	}{
		{
			Name:   "success",
			UserId: 1,
			User: user.User{
				Id:       1,
				Email:    "test@gmail.com",
				Password: "jdfbjfbj",
			},
			IsValid: true,
		},
		{
			Name:    "failed",
			UserId:  4,
			User:    user.User{},
			IsValid: false,
			Error:   errors.New("failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)
			interactor := user.NewInteractor(repo)

			repo.On("FindById", tc.UserId).Return(tc.User, tc.Error)
			res, err := interactor.ValidateUser(tc.UserId)
			assert.Equal(t, tc.Error, err)
			assert.Equal(t, tc.IsValid, res)
		})
	}
}

func TestGetUserById(t *testing.T) {
	testCases := []struct {
		Name     string
		UserId   int
		User     user.User
		HttpCode int
		Error    error
	}{
		{
			Name:   "success",
			UserId: 2,
			User: user.User{
				Id:       2,
				Email:    "test@gmail.com",
				Password: "fjbdjfb232j",
			},
			HttpCode: http.StatusOK,
		},
		{
			Name:     "failed",
			UserId:   1,
			User:     user.User{},
			HttpCode: http.StatusInternalServerError,
			Error:    errors.New("failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)
			interactor := user.NewInteractor(repo)

			repo.On("FindById", tc.UserId).Return(tc.User, tc.Error)
			res, httpCode, err := interactor.GetUserById(tc.UserId)
			assert.Equal(t, tc.Error, err)
			assert.Equal(t, tc.HttpCode, httpCode)
			assert.Equal(t, tc.User, res)
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	testCases := []struct {
		Name     string
		Email    string
		User     user.User
		HttpCode int
		Error    error
	}{
		{
			Name:  "success",
			Email: "test@gmail.com",
			User: user.User{
				Id:       1,
				Email:    "test@gmail.com",
				Password: "fnjdfbdj",
			},
			HttpCode: http.StatusOK,
		},
		{
			Name:     "failed",
			Email:    "another@gmail.com",
			User:     user.User{},
			HttpCode: http.StatusInternalServerError,
			Error:    errors.New("failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)
			interactor := user.NewInteractor(repo)

			repo.On("FindByEmail", tc.Email).Return(tc.User, tc.Error)
			res, httpCode, err := interactor.GetUserByEmail(tc.Email)
			assert.Equal(t, tc.Error, err)
			assert.Equal(t, tc.HttpCode, httpCode)
			assert.Equal(t, tc.User, res)
		})
	}
}
