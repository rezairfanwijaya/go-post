package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-post/internal/auth"
	"go-post/internal/handler"
	"go-post/internal/post"
	postmock "go-post/internal/post/mocks"
	"go-post/internal/user"
	usermock "go-post/internal/user/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestSignUpSuccess(t *testing.T) {
	pRepo := postmock.NewPostRepository(t)
	pInteractor := post.NewInteractor(pRepo)
	uInteractor := usermock.NewInteractor(t)
	auth := auth.NewAuthService()
	h := handler.NewUserHandler(pInteractor, uInteractor, auth)

	testCase := struct {
		Name         string
		Input        user.InputUserSignUp
		ExpectedCode int
	}{
		Name: "success",
		Input: user.InputUserSignUp{
			Email:    "test@gmail.com",
			Password: "123456789",
		},
		ExpectedCode: http.StatusOK,
	}

	t.Run(testCase.Name, func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		inputJson, err := json.Marshal(testCase.Input)
		assert.NoError(t, err)
		uInteractor.On("CreateUser", mock.AnythingOfType("User")).Return(http.StatusOK, nil)
		req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer(inputJson))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.POST("/users/signup", h.SignUp)
		r.ServeHTTP(rec, req)
		assert.Equal(t, testCase.ExpectedCode, rec.Code)
	})
}

func TestSignUpFailed(t *testing.T) {
	pRepo := postmock.NewPostRepository(t)
	pInteractor := post.NewInteractor(pRepo)
	uInteractor := usermock.NewInteractor(t)
	auth := auth.NewAuthService()
	h := handler.NewUserHandler(pInteractor, uInteractor, auth)

	testCase := struct {
		Name         string
		Input        user.InputUserSignUp
		ExpectedCode int
	}{
		Name: "failed",
		Input: user.InputUserSignUp{
			Email: "test@gmail.com",
		},
		ExpectedCode: http.StatusInternalServerError,
	}

	t.Run(testCase.Name, func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		inputJson, err := json.Marshal(testCase.Input)
		assert.NoError(t, err)
		uInteractor.On("CreateUser", mock.AnythingOfType("User")).Return(http.StatusInternalServerError, errors.New("failed"))
		req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer(inputJson))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.POST("/users/signup", h.SignUp)
		r.ServeHTTP(rec, req)
		assert.Equal(t, testCase.ExpectedCode, rec.Code)
	})
}

func TestLoginSuccess(t *testing.T) {
	pRepo := postmock.NewPostRepository(t)
	pInteractor := post.NewInteractor(pRepo)
	uInteractor := usermock.NewInteractor(t)
	auth := auth.NewAuthService()
	h := handler.NewUserHandler(pInteractor, uInteractor, auth)

	testCase := struct {
		Name         string
		Input        user.InputUserLogin
		ExpectedCode int
	}{
		Name: "success",
		Input: user.InputUserLogin{
			Email:    "test@gmail.com",
			Password: "12345678",
		},
		ExpectedCode: http.StatusOK,
	}

	bytePass, _ := bcrypt.GenerateFromPassword([]byte(testCase.Input.Password), 10)

	u := user.User{
		Id:       1,
		Email:    "test@gmail.com",
		Password: string(bytePass),
	}

	t.Run(testCase.Name, func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		inputJson, err := json.Marshal(testCase.Input)
		assert.NoError(t, err)

		uInteractor.On("GetUserByEmail", testCase.Input.Email).Return(u, http.StatusOK, nil)
		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(inputJson))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.POST("/users/login", h.Login)
		r.ServeHTTP(rec, req)
		assert.Equal(t, testCase.ExpectedCode, rec.Code)

	})
}

func TestLoginFailed(t *testing.T) {
	pRepo := postmock.NewPostRepository(t)
	pInteractor := post.NewInteractor(pRepo)
	uInteractor := usermock.NewInteractor(t)
	auth := auth.NewAuthService()
	h := handler.NewUserHandler(pInteractor, uInteractor, auth)

	testCase := struct {
		Name         string
		Input        user.InputUserLogin
		ExpectedCode int
	}{
		Name: "success",
		Input: user.InputUserLogin{
			Email:    "test@gmail.com",
			Password: "12345678",
		},
		ExpectedCode: http.StatusBadRequest,
	}

	u := user.User{
		Id:       1,
		Email:    "test@gmail.com",
		Password: "12345678",
	}

	t.Run(testCase.Name, func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		inputJson, err := json.Marshal(testCase.Input)
		assert.NoError(t, err)

		uInteractor.On("GetUserByEmail", testCase.Input.Email).Return(u, http.StatusOK, nil)
		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(inputJson))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.POST("/users/login", h.Login)
		r.ServeHTTP(rec, req)
		assert.Equal(t, testCase.ExpectedCode, rec.Code)

	})
}
