package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-post/internal/handler"
	"go-post/internal/post"
	postmocks "go-post/internal/post/mocks"
	"go-post/internal/user"
	usermocks "go-post/internal/user/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CreatePost_Success(t *testing.T) {
	uInteractor := usermocks.NewInteractor(t)
	pInteractor := postmocks.NewInteractor(t)

	h := handler.NewPostHandler(uInteractor, pInteractor)

	type expectation struct {
		ResponseBody []byte
		HttpCode     int
	}

	testCase := struct {
		Name   string
		Input  post.InputCreatePost
		Result expectation
	}{
		Name: "success",
		Input: post.InputCreatePost{
			Title:   "test",
			Content: "test",
		},
		Result: expectation{
			ResponseBody: []byte(`{"meta":{"code":200,"status":"success"},"data":"success"}`),
			HttpCode:     http.StatusOK,
		},
	}

	t.Run(testCase.Name, func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", 1)
		})

		inputJson, err := json.Marshal(testCase.Input)
		assert.NoError(t, err)

		pInteractor.On("CreatePost", mock.AnythingOfType("Post")).Return(http.StatusOK, nil)
		req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(inputJson))
		req.Header.Add("Content-Type", "application/post")
		rec := httptest.NewRecorder()

		r.POST("/posts", h.CreatePost)
		r.ServeHTTP(rec, req)

		assert.Equal(t, bytes.NewBuffer(testCase.Result.ResponseBody), rec.Body)
		assert.Equal(t, testCase.Result.HttpCode, rec.Code)
	})
}

func Test_CereatePost_Failed(t *testing.T) {
	uInteractor := usermocks.NewInteractor(t)
	pInteractor := postmocks.NewInteractor(t)
	h := handler.NewPostHandler(uInteractor, pInteractor)

	type expecation struct {
		HttpMethod   int
		ResponseBody []byte
	}

	testCase := struct {
		Name   string
		Input  post.InputCreatePost
		Result expecation
	}{
		Name: "failed",
		Input: post.InputCreatePost{
			Title:   "test",
			Content: "test",
		},
		Result: expecation{
			HttpMethod:   http.StatusInternalServerError,
			ResponseBody: []byte(`{"meta":{"code":500,"status":"error"},"data":"failed"}`),
		},
	}

	t.Run(testCase.Name, func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", 1)
		})

		jsonInput, err := json.Marshal(testCase.Input)
		assert.Nil(t, err)

		pInteractor.On("CreatePost", mock.AnythingOfType("Post")).Return(http.StatusInternalServerError, errors.New("failed"))
		req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(jsonInput))
		rec := httptest.NewRecorder()

		r.POST("/posts", h.CreatePost)
		r.ServeHTTP(rec, req)

		assert.Equal(t, bytes.NewBuffer(testCase.Result.ResponseBody), rec.Body)
		assert.Equal(t, testCase.Result.HttpMethod, rec.Code)
	})
}

func Test_GetPost_Success(t *testing.T) {
	pInteractor := postmocks.NewInteractor(t)
	uInteractor := usermocks.NewInteractor(t)
	h := handler.NewPostHandler(uInteractor, pInteractor)

	post := post.Post{
		Id:      1,
		UserId:  1,
		Title:   "test",
		Content: "test",
	}

	user := user.User{
		Id:       1,
		Email:    "test@gmail.com",
		Password: "123456789",
	}

	responseBody := []byte(`{"meta":{"code":200,"status":"success"},"data":{"post":{"id":1,"title":"test","content":"test"},"user":{"id":1,"email":"test@gmail.com"}}}`)

	t.Run("success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", 1)
		})

		pInteractor.On("GetPost", 1, 1).Return(post, http.StatusOK, nil)
		uInteractor.On("GetUserById", 1).Return(user, http.StatusOK, nil)

		uri := "/posts/1"
		req := httptest.NewRequest(http.MethodGet, uri, nil)
		rec := httptest.NewRecorder()

		r.GET("/posts/:id", h.GetPost)
		r.ServeHTTP(rec, req)

		assert.Equal(t, bytes.NewBuffer(responseBody), rec.Body)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func Test_GetPost_Err_PostId(t *testing.T) {
	pInteractor := postmocks.NewInteractor(t)
	uInteractor := usermocks.NewInteractor(t)
	h := handler.NewPostHandler(uInteractor, pInteractor)

	responseBody := []byte(`{"meta":{"code":400,"status":"error in convert id"},"data":"strconv.Atoi: parsing \"bnm\": invalid syntax"}`)

	t.Run("error_post_id", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", 1)
		})

		uri := "/posts/bnm"
		req := httptest.NewRequest(http.MethodGet, uri, nil)
		rec := httptest.NewRecorder()

		r.GET("/posts/:id", h.GetPost)
		r.ServeHTTP(rec, req)

		assert.Equal(t, bytes.NewBuffer(responseBody), rec.Body)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func Test_GetPost_PostFailed(t *testing.T) {
	pInteractor := postmocks.NewInteractor(t)
	uInteractor := usermocks.NewInteractor(t)
	h := handler.NewPostHandler(uInteractor, pInteractor)

	expectation := []byte(`{"meta":{"code":500,"status":"error"},"data":"failed"}`)

	t.Run("error_post_not_found", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", 1)
		})

		pInteractor.On("GetPost", 1, 1).Return(post.Post{}, http.StatusInternalServerError, errors.New("failed"))

		uri := "/posts/1"
		req := httptest.NewRequest(http.MethodGet, uri, nil)
		rec := httptest.NewRecorder()

		r.GET("/posts/:id", h.GetPost)
		r.ServeHTTP(rec, req)

		assert.Equal(t, bytes.NewBuffer(expectation), rec.Body)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func Test_GetPost_UserFailed(t *testing.T) {
	pInteractor := postmocks.NewInteractor(t)
	uInteractor := usermocks.NewInteractor(t)
	h := handler.NewPostHandler(uInteractor, pInteractor)

	post := post.Post{
		Id:      1,
		UserId:  1,
		Title:   "test",
		Content: "test",
	}

	expectation := []byte(`{"meta":{"code":500,"status":"error"},"data":"failed"}`)

	t.Run("error_post_not_found", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", 1)
		})

		pInteractor.On("GetPost", 1, 1).Return(post, http.StatusOK, nil)
		uInteractor.On("GetUserById", 1).Return(user.User{}, http.StatusInternalServerError, errors.New("failed"))

		uri := "/posts/1"
		req := httptest.NewRequest(http.MethodGet, uri, nil)
		rec := httptest.NewRecorder()

		r.GET("/posts/:id", h.GetPost)
		r.ServeHTTP(rec, req)

		assert.Equal(t, bytes.NewBuffer(expectation), rec.Body)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func Test_DeletePost_Success(t *testing.T) {
	uInteractor := usermocks.NewInteractor(t)
	pInteractor := postmocks.NewInteractor(t)
	h := handler.NewPostHandler(uInteractor, pInteractor)

	expectation := []byte(`{"meta":{"code":200,"status":"success"},"data":"success"}`)

	t.Run("success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", 1)
		})

		pInteractor.On("DeletePost", 1, 1).Return(http.StatusOK, nil)

		uri := "/posts/1"
		req := httptest.NewRequest(http.MethodDelete, uri, nil)
		rec := httptest.NewRecorder()

		r.DELETE("/posts/:id", h.DeletePost)
		r.ServeHTTP(rec, req)

		assert.Equal(t, bytes.NewBuffer(expectation), rec.Body)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func Test_DeletePost_Failed(t *testing.T) {
	pInteractor := postmocks.NewInteractor(t)
	uInteractor := usermocks.NewInteractor(t)
	h := handler.NewPostHandler(uInteractor, pInteractor)

	expectation := []byte(`{"meta":{"code":500,"status":"error"},"data":"failed"}`)

	t.Run("failed", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", 1)
		})

		t.Run("failed", func(t *testing.T) {
			pInteractor.On("DeletePost", 1, 1).Return(http.StatusInternalServerError, errors.New("failed"))

			uri := "/posts/1"
			req := httptest.NewRequest(http.MethodDelete, uri, nil)
			rec := httptest.NewRecorder()

			r.DELETE("/posts/:id", h.DeletePost)
			r.ServeHTTP(rec, req)

			assert.Equal(t, bytes.NewBuffer(expectation), rec.Body)
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		})
	})
}

func Test_DeletePost_Err_PostId(t *testing.T) {
	pInteractor := postmocks.NewInteractor(t)
	uInteractor := usermocks.NewInteractor(t)
	h := handler.NewPostHandler(uInteractor, pInteractor)

	expectation := []byte(`{"meta":{"code":400,"status":"error in convert id"},"data":"strconv.Atoi: parsing \"jxc\": invalid syntax"}`)

	t.Run("err_post_id", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", 1)
		})

		uri := "/posts/jxc"
		req := httptest.NewRequest(http.MethodDelete, uri, nil)
		rec := httptest.NewRecorder()

		r.DELETE("/posts/:id", h.DeletePost)
		r.ServeHTTP(rec, req)

		assert.Equal(t, bytes.NewBuffer(expectation), rec.Body)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func Test_UpdatePost_Success(t *testing.T) {
	pInteractor := postmocks.NewInteractor(t)
	uInteractor := usermocks.NewInteractor(t)
	h := handler.NewPostHandler(uInteractor, pInteractor)

	input := post.InputUpdatePost{
		UserId:  1,
		Title:   "test update",
		Content: "test update",
	}

	postById := post.Post{
		Id:      1,
		UserId:  1,
		Title:   "test",
		Content: "test",
	}

	postUpdated := post.Post{
		UserId:  1,
		Title:   "test update",
		Content: "test update",
	}

	expectation := []byte(`{"meta":{"code":200,"status":"success"},"data":{"id":0,"title":"test update","content":"test update"}}`)

	t.Run("success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", 1)
		})

		pInteractor.On("GetPost", 1, 1).Return(postById, http.StatusOK, nil)
		pInteractor.On("UpdatePost", 1, 1, mock.AnythingOfType("Post")).Return(postUpdated, http.StatusOK, nil)

		uri := "/posts/1"
		inputJson, err := json.Marshal(input)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, uri, bytes.NewBuffer([]byte(inputJson)))
		rec := httptest.NewRecorder()

		r.PUT("/posts/:id", h.UpdatePost)
		r.ServeHTTP(rec, req)

		assert.Equal(t, bytes.NewBuffer(expectation), rec.Body)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func Test_UpdatePost_Err_PostId(t *testing.T) {
	pInteractor := postmocks.NewInteractor(t)
	uInteractor := usermocks.NewInteractor(t)
	h := handler.NewPostHandler(uInteractor, pInteractor)

	input := post.InputUpdatePost{
		UserId:  1,
		Title:   "test update",
		Content: "test update",
	}

	expectation := []byte(`{"meta":{"code":400,"status":"error in convert id"},"data":"strconv.Atoi: parsing \"xcv\": invalid syntax"}`)

	t.Run("err_post_id", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", 1)
		})

		uri := "/posts/xcv"
		inputJson, err := json.Marshal(input)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPut, uri, bytes.NewBuffer(inputJson))
		rec := httptest.NewRecorder()

		r.PUT("/posts/:id", h.UpdatePost)
		r.ServeHTTP(rec, req)

		assert.Equal(t, bytes.NewBuffer(expectation), rec.Body)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func Test_UpdatePost_GetPostFailed(t *testing.T) {
	pInteractor := postmocks.NewInteractor(t)
	uInteractor := usermocks.NewInteractor(t)
	h := handler.NewPostHandler(uInteractor, pInteractor)

	input := post.InputUpdatePost{
		UserId:  1,
		Title:   "test update",
		Content: "test update",
	}

	expectation := []byte(`{"meta":{"code":500,"status":"error"},"data":"failed"}`)

	t.Run("update_getpost_failed", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", 1)
		})

		uri := "/posts/1"
		inputJson, err := json.Marshal(input)
		assert.NoError(t, err)

		pInteractor.On("GetPost", 1, 1).Return(post.Post{}, http.StatusInternalServerError, errors.New("failed"))

		req := httptest.NewRequest(http.MethodPut, uri, bytes.NewBuffer(inputJson))
		rec := httptest.NewRecorder()

		r.PUT("/posts/:id", h.UpdatePost)
		r.ServeHTTP(rec, req)

		assert.Equal(t, bytes.NewBuffer(expectation), rec.Body)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func Test_UpdatePost_Failed(t *testing.T) {
	pInteractor := postmocks.NewInteractor(t)
	uInteractor := usermocks.NewInteractor(t)
	h := handler.NewPostHandler(uInteractor, pInteractor)

	input := post.InputUpdatePost{
		UserId:  1,
		Title:   "test update",
		Content: "test update",
	}

	postById := post.Post{
		Id:      1,
		UserId:  1,
		Title:   "test",
		Content: "test",
	}

	expectation := []byte(`{"meta":{"code":500,"status":"error"},"data":"failed"}`)

	t.Run("failed", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", 1)
		})

		uri := "/posts/1"
		inputJson, err := json.Marshal(input)
		assert.NoError(t, err)

		pInteractor.On("GetPost", 1, 1).Return(postById, http.StatusOK, nil)
		pInteractor.On("UpdatePost", 1, 1, mock.AnythingOfType("Post")).Return(post.Post{}, http.StatusInternalServerError, errors.New("failed"))

		req := httptest.NewRequest(http.MethodPut, uri, bytes.NewBuffer(inputJson))
		rec := httptest.NewRecorder()

		r.PUT("/posts/:id", h.UpdatePost)
		r.ServeHTTP(rec, req)

		assert.Equal(t, bytes.NewBuffer(expectation), rec.Body)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
