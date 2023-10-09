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
	"log"
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
		Name: "failed",
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

		uInteractor.On("GetUserByEmail", testCase.Input.Email).Return(u, http.StatusBadRequest, errors.New("failed"))
		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(inputJson))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.POST("/users/login", h.Login)
		r.ServeHTTP(rec, req)
		assert.Equal(t, testCase.ExpectedCode, rec.Code)
	})
}

func TestGetUserWithPostsSuccess(t *testing.T) {
	pInteractor := postmock.NewInteractor(t)
	uInteractor := usermock.NewInteractor(t)
	auth := auth.NewAuthService()
	h := handler.NewUserHandler(pInteractor, uInteractor, auth)

	testCase := struct {
		Name         string
		UserId       int
		ExpectedCode int
	}{
		Name:         "success",
		UserId:       1,
		ExpectedCode: http.StatusOK,
	}

	posts := []post.Post{
		{
			Id:      1,
			UserId:  1,
			Title:   "test 1",
			Content: "test 1",
		},
		{
			Id:      2,
			UserId:  1,
			Title:   "test 2",
			Content: "test 2",
		},
	}

	user := user.User{
		Id:       1,
		Email:    "test@gmail.com",
		Password: "123456789",
	}

	t.Run(testCase.Name, func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		pInteractor.On("GetPostByUserId", testCase.UserId).Return(posts, http.StatusOK, nil)
		uInteractor.On("GetUserById", testCase.UserId).Return(user, http.StatusOK, nil)

		req := httptest.NewRequest(http.MethodGet, "/users/posts", nil)
		rec := httptest.NewRecorder()

		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", 1)
		})

		r.GET("/users/posts", h.GetUserWithPosts)

		r.ServeHTTP(rec, req)
	})
}

func TestGetUserWithPostsFailedFindUser(t *testing.T) {
	pInteractor := postmock.NewInteractor(t)
	uInteractor := usermock.NewInteractor(t)
	auth := auth.NewAuthService()
	h := handler.NewUserHandler(pInteractor, uInteractor, auth)

	testCase := struct {
		Name         string
		UserId       int
		ExpectedCode int
	}{
		Name:         "failed get user",
		UserId:       1,
		ExpectedCode: http.StatusOK,
	}

	posts := []post.Post{
		{
			Id:      1,
			UserId:  1,
			Title:   "test 1",
			Content: "test 1",
		},
		{
			Id:      2,
			UserId:  1,
			Title:   "test 2",
			Content: "test 2",
		},
	}

	t.Run(testCase.Name, func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		pInteractor.On("GetPostByUserId", testCase.UserId).Return(posts, http.StatusOK, nil)
		uInteractor.On("GetUserById", testCase.UserId).Return(user.User{}, http.StatusInternalServerError, errors.New("failed"))

		req := httptest.NewRequest(http.MethodGet, "/users/posts", nil)
		rec := httptest.NewRecorder()

		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", 1)
		})

		r.GET("/users/posts", h.GetUserWithPosts)
		log.Println(rec)
		r.ServeHTTP(rec, req)
	})
}

func TestGetUserWithPostsFailedFindPosts(t *testing.T) {
	pInteractor := postmock.NewInteractor(t)
	uInteractor := usermock.NewInteractor(t)
	auth := auth.NewAuthService()
	h := handler.NewUserHandler(pInteractor, uInteractor, auth)

	testCase := struct {
		Name         string
		UserId       int
		ExpectedCode int
	}{
		Name:         "failed get posts",
		UserId:       1,
		ExpectedCode: http.StatusOK,
	}

	t.Run(testCase.Name, func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		pInteractor.On("GetPostByUserId", testCase.UserId).Return([]post.Post{}, http.StatusInternalServerError, errors.New("failed"))

		req := httptest.NewRequest(http.MethodGet, "/users/posts", nil)
		rec := httptest.NewRecorder()

		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", 1)
		})

		r.GET("/users/posts", h.GetUserWithPosts)
		r.ServeHTTP(rec, req)
	})
}
