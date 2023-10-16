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
	repo := mocks.NewUserRepository(t)
	interactor := user.NewInteractor(repo)

	testCases := []struct {
		Name     string
		User     user.User
		HttpCode int
		WantErr  bool
	}{
		{
			Name: "success",
			User: user.User{
				Email:    "test@gmail.com",
				Password: "djbgj121-232j",
			},
			HttpCode: http.StatusOK,
			WantErr:  false,
		},
		{
			Name:     "failed",
			User:     user.User{},
			HttpCode: http.StatusInternalServerError,
			WantErr:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if !testCase.WantErr {
				repo.On("Save", testCase.User).Return(testCase.User, nil)
				_, httpCode, err := interactor.CreateUser(testCase.User)
				assert.Nil(t, err)
				assert.Equal(t, testCase.HttpCode, httpCode)
			} else {
				repo.On("Save", testCase.User).Return(testCase.User, errors.New("failed"))
				_, httpCode, err := interactor.CreateUser(testCase.User)
				assert.NotNil(t, err)
				assert.Equal(t, testCase.HttpCode, httpCode)
			}
		})
	}
}

func TestValidateUser(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	interactor := user.NewInteractor(repo)

	testCases := []struct {
		Name    string
		UserId  int
		User    user.User
		IsValid bool
		WantErr bool
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
			WantErr: false,
		},
		{
			Name:    "failed",
			UserId:  4,
			User:    user.User{},
			IsValid: false,
			WantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if !testCase.WantErr {
				repo.On("FindById", testCase.UserId).Return(testCase.User, nil)
				isValid, err := interactor.ValidateUser(testCase.UserId)
				assert.Equal(t, testCase.IsValid, isValid)
				assert.Nil(t, err)
			} else {
				repo.On("FindById", testCase.UserId).Return(testCase.User, errors.New("failed"))
				isValid, err := interactor.ValidateUser(testCase.UserId)
				assert.Equal(t, isValid, testCase.IsValid)
				assert.NotNil(t, err)
			}
		})
	}
}

func TestGetUserById(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	interactor := user.NewInteractor(repo)

	testCases := []struct {
		Name     string
		UserId   int
		User     user.User
		HttpCode int
		WantErr  bool
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
			WantErr:  false,
		},
		{
			Name:     "failed",
			UserId:   1,
			User:     user.User{},
			HttpCode: http.StatusInternalServerError,
			WantErr:  true,
		},
	}

	for _, testCase := range testCases {
		if !testCase.WantErr {
			repo.On("FindById", testCase.UserId).Return(testCase.User, nil)
			user, httpCode, err := interactor.GetUserById(testCase.UserId)
			assert.Equal(t, testCase.User, user)
			assert.Equal(t, testCase.HttpCode, httpCode)
			assert.Nil(t, err)
		} else {
			repo.On("FindById", testCase.UserId).Return(testCase.User, errors.New("failed"))
			user, httpCode, err := interactor.GetUserById(testCase.UserId)
			assert.Equal(t, testCase.User, user)
			assert.Equal(t, testCase.HttpCode, httpCode)
			assert.NotNil(t, err)
		}
	}
}

func TestGetUserByEmail(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	interactor := user.NewInteractor(repo)

	testCases := []struct {
		Name     string
		Email    string
		User     user.User
		HttpCode int
		WantErr  bool
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
			WantErr:  false,
		},
		{
			Name:     "failed",
			Email:    "another@gmail.com",
			User:     user.User{},
			HttpCode: http.StatusInternalServerError,
			WantErr:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if !testCase.WantErr {
				repo.On("FindByEmail", testCase.Email).Return(testCase.User, nil)
				user, httpCode, err := interactor.GetUserByEmail(testCase.Email)
				assert.Equal(t, testCase.User, user)
				assert.Equal(t, testCase.HttpCode, httpCode)
				assert.Nil(t, err)
			} else {
				repo.On("FindByEmail", testCase.Email).Return(testCase.User, errors.New("failed"))
				user, httpCode, err := interactor.GetUserByEmail(testCase.Email)
				assert.Equal(t, testCase.User, user)
				assert.Equal(t, testCase.HttpCode, httpCode)
				assert.NotNil(t, err)
			}
		})
	}
}
