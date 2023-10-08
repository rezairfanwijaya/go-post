package post_test

import (
	"errors"
	"go-post/internal/post"
	"go-post/internal/post/mocks"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUser(t *testing.T) {
	repo := mocks.NewPostRepository(t)
	interactor := post.NewInteractor(repo)

	post := post.Post{
		Id:      1,
		UserId:  1,
		Title:   "test",
		Content: "test",
	}

	isValid := interactor.ValidateUser(1, post)
	assert.Equal(t, true, isValid)
}

func TestCreatePost(t *testing.T) {
	repo := mocks.NewPostRepository(t)
	interactor := post.NewInteractor(repo)

	testCases := []struct {
		Name     string
		post     post.Post
		HttpCode int
		WantErr  bool
	}{
		{
			Name: "success",
			post: post.Post{
				Id:      1,
				UserId:  2,
				Title:   "test",
				Content: "test",
			},
			HttpCode: http.StatusOK,
			WantErr:  false,
		},
		{
			Name:     "failed",
			post:     post.Post{},
			HttpCode: http.StatusInternalServerError,
			WantErr:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if !testCase.WantErr {
				repo.On("Save", testCase.post).Return(nil)
				httpCode, err := interactor.CreatePost(testCase.post)
				assert.Equal(t, testCase.HttpCode, httpCode)
				assert.Nil(t, err)
			} else {
				repo.On("Save", testCase.post).Return(errors.New("failed"))
				httpCode, err := interactor.CreatePost(testCase.post)
				assert.Equal(t, testCase.HttpCode, httpCode)
				assert.NotNil(t, err)
			}
		})
	}
}

func TestGetPost(t *testing.T) {
	repo := mocks.NewPostRepository(t)
	interactor := post.NewInteractor(repo)

	type params struct {
		PostId int
		UserId int
	}

	type response struct {
		Post     post.Post
		HttpCode int
	}

	testCases := []struct {
		Name            string
		Params          params
		Response        response
		WantErr         bool
		IsAunauthorized bool
	}{
		{
			Name: "success",
			Params: params{
				PostId: 1,
				UserId: 1,
			},
			Response: response{
				Post: post.Post{
					Id:      1,
					UserId:  1,
					Title:   "test",
					Content: "test",
				},
				HttpCode: http.StatusOK,
			},
			WantErr:         false,
			IsAunauthorized: true,
		},
		{
			Name: "unauthorized",
			Params: params{
				PostId: 1,
				UserId: 2,
			},
			Response: response{
				Post:     post.Post{},
				HttpCode: http.StatusUnauthorized,
			},
			WantErr:         true,
			IsAunauthorized: false,
		},
		{
			Name:   "failed",
			Params: params{},
			Response: response{
				Post:     post.Post{},
				HttpCode: http.StatusInternalServerError,
			},
			WantErr:         true,
			IsAunauthorized: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if !testCase.WantErr {
				repo.On("FindByPostId", testCase.Params.PostId).Return(testCase.Response.Post, nil)
				post, httpCode, err := interactor.GetPost(testCase.Params.UserId, testCase.Params.PostId)
				assert.Equal(t, testCase.Response.Post, post)
				assert.Equal(t, testCase.Response.HttpCode, httpCode)
				assert.Nil(t, err)
			} else if !testCase.IsAunauthorized {
				repo.On("FindByPostId", testCase.Params.PostId).Return(testCase.Response.Post, nil)
				post, httpCode, err := interactor.GetPost(testCase.Params.UserId, testCase.Params.PostId)
				assert.Equal(t, testCase.Response.Post, post)
				assert.Equal(t, testCase.Response.HttpCode, httpCode)
				assert.NotNil(t, err)
			} else {
				repo.On("FindByPostId", testCase.Params.PostId).Return(testCase.Response.Post, errors.New("failed"))
				post, httpCode, err := interactor.GetPost(testCase.Params.UserId, testCase.Params.PostId)
				assert.Equal(t, testCase.Response.Post, post)
				assert.Equal(t, testCase.Response.HttpCode, httpCode)
				assert.NotNil(t, err)
			}
		})
	}
}

func TestGetPostByUserId(t *testing.T) {
	repo := mocks.NewPostRepository(t)
	interactor := post.NewInteractor(repo)

	testCases := []struct {
		Name     string
		UserId   int
		Posts    []post.Post
		HttpCode int
		WantErr  bool
	}{
		{
			Name:   "success",
			UserId: 1,
			Posts: []post.Post{
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
			},
			HttpCode: http.StatusOK,
			WantErr:  false,
		},
		{
			Name:     "failed",
			UserId:   8,
			Posts:    []post.Post{},
			HttpCode: http.StatusInternalServerError,
			WantErr:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if !testCase.WantErr {
				repo.On("FindByUserId", testCase.UserId).Return(testCase.Posts, nil)
				posts, httpCode, err := interactor.GetPostByUserId(testCase.UserId)
				assert.Equal(t, testCase.Posts, posts)
				assert.Equal(t, testCase.HttpCode, httpCode)
				assert.Nil(t, err)
			} else {
				repo.On("FindByUserId", testCase.UserId).Return(testCase.Posts, errors.New("failed"))
				posts, httpCode, err := interactor.GetPostByUserId(testCase.UserId)
				assert.Equal(t, testCase.Posts, posts)
				assert.Equal(t, testCase.HttpCode, httpCode)
				assert.NotNil(t, err)
			}
		})
	}
}

func TestUpdatePost(t *testing.T) {
	repo := mocks.NewPostRepository(t)
	interactor := post.NewInteractor(repo)

	type params struct {
		UserId int
		PostId int
		Post   post.Post
	}

	type response struct {
		Post     post.Post
		HttpCode int
	}

	testCases := []struct {
		Name           string
		Params         params
		Response       response
		WantErr        bool
		IsUnauthorized bool
	}{
		{
			Name: "success",
			Params: params{
				UserId: 1,
				PostId: 1,
				Post: post.Post{
					Id:      1,
					UserId:  1,
					Title:   "test",
					Content: "test",
				},
			},
			Response: response{
				Post: post.Post{
					Id:      1,
					UserId:  1,
					Title:   "test",
					Content: "test",
				},
				HttpCode: http.StatusOK,
			},
			WantErr:        false,
			IsUnauthorized: true,
		},
		{
			Name: "unauthorized",
			Params: params{
				UserId: 2,
				PostId: 1,
				Post: post.Post{
					Id:      1,
					UserId:  8,
					Title:   "test",
					Content: "test",
				},
			},
			Response: response{
				Post:     post.Post{},
				HttpCode: http.StatusUnauthorized,
			},
			WantErr:        true,
			IsUnauthorized: false,
		},
		{
			Name:   "failed",
			Params: params{},
			Response: response{
				Post:     post.Post{},
				HttpCode: http.StatusInternalServerError,
			},
			WantErr:        true,
			IsUnauthorized: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if !testCase.WantErr {
				repo.On("Update", testCase.Params.PostId, testCase.Params.Post).Return(nil)

				post, httpCode, err := interactor.UpdatePost(testCase.Params.PostId, testCase.Params.UserId, testCase.Params.Post)
				assert.Equal(t, testCase.Response.Post, post)
				assert.Equal(t, testCase.Response.HttpCode, httpCode)
				assert.Nil(t, err)
			} else if !testCase.IsUnauthorized {
				post, httpCode, err := interactor.UpdatePost(testCase.Params.PostId, testCase.Params.UserId, testCase.Params.Post)
				assert.Equal(t, testCase.Response.Post, post)
				assert.Equal(t, testCase.Response.HttpCode, httpCode)
				assert.NotNil(t, err)
			} else {
				repo.On("Update", testCase.Params.PostId, testCase.Params.Post).Return(errors.New("failed"))

				post, httpCode, err := interactor.UpdatePost(testCase.Params.PostId, testCase.Params.UserId, testCase.Params.Post)
				assert.Equal(t, testCase.Response.Post, post)
				assert.Equal(t, testCase.Response.HttpCode, httpCode)
				assert.NotNil(t, err)
			}
		})
	}
}

func TestDeletePost(t *testing.T) {
	repo := mocks.NewPostRepository(t)
	interactor := post.NewInteractor(repo)

	type Params struct {
		PostId int
		UserId int
	}

	testCases := []struct {
		Name           string
		Params         Params
		HttpCode       int
		IsUnauthorized bool
		WantErr        bool
	}{
		{
			Name: "success",
			Params: Params{
				PostId: 1,
				UserId: 2,
			},
			HttpCode:       http.StatusOK,
			IsUnauthorized: true,
			WantErr:        false,
		},
		{
			Name: "unauthorized",
			Params: Params{
				PostId: 1,
				UserId: 78,
			},
			HttpCode:       http.StatusUnauthorized,
			IsUnauthorized: false,
			WantErr:        true,
		},
		{
			Name: "failed",
			Params: Params{
				PostId: 1,
				UserId: 2,
			},
			HttpCode:       http.StatusInternalServerError,
			IsUnauthorized: true,
			WantErr:        true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if !testCase.WantErr {
				repo.On("FindByPostId", testCase.Params.PostId).Return(post.Post{Id: 1, UserId: 2}, nil)
				repo.On("Delete", testCase.Params.PostId).Return(nil)

				httpCode, err := interactor.DeletePost(testCase.Params.PostId, testCase.Params.UserId)
				assert.Equal(t, testCase.HttpCode, httpCode)
				assert.Nil(t, err)
			} else if !testCase.IsUnauthorized {
				repo.On("FindByPostId", testCase.Params.PostId).Return(post.Post{Id: 1, UserId: 3}, nil)

				httpCode, err := interactor.DeletePost(testCase.Params.PostId, testCase.Params.UserId)
				assert.Equal(t, testCase.HttpCode, httpCode)
				assert.NotNil(t, err)
			} else {
				repo.On("FindByPostId", testCase.Params.PostId).Return(post.Post{Id: 1, UserId: 2}, nil)
				
				repo.On("Delete", testCase.Params.PostId).Return(errors.New("failed"))
				httpCode, err := interactor.DeletePost(testCase.Params.PostId, testCase.Params.UserId)
				t.Log(httpCode, err)
			}
		})
	}
}
