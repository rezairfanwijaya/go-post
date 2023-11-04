package post_test

import (
	"go-post/internal/post"
	"go-post/internal/post/mocks"
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
	testCases := []struct {
		Name        string
		Post        post.Post
		ExpectedErr error
		ExpecredRes post.Post
	}{
		{
			Name: "success",
			Post: post.Post{
				Id:      1,
				UserId:  2,
				Title:   "test",
				Content: "test",
			},
			ExpecredRes: post.Post{
				Id:      1,
				UserId:  2,
				Title:   "test",
				Content: "test",
			},
		},
		{
			Name:        "failed",
			Post:        post.Post{},
			ExpecredRes: post.Post{},
			ExpectedErr: post.ErrDatabaseFailure,
		},
	}

	for _, tc := range testCases {
		r := mocks.NewPostRepository(t)
		i := post.NewInteractor(r)

		r.On("Save", tc.Post).Return(tc.ExpecredRes, tc.ExpectedErr)
		res, err := i.CreatePost(tc.Post)
		assert.Equal(t, tc.ExpectedErr, err)
		assert.Equal(t, tc.ExpecredRes, res)
	}
}

func TestGetPost(t *testing.T) {
	p := post.Post{
		Id:      1,
		UserId:  1,
		Title:   "test",
		Content: "test",
	}
	testCases := []struct {
		Name         string
		PostId       int
		UserId       int
		ExpectedRes  post.Post
		ExpectedPost post.Post
		ExpectedErr  error
		ErrRepoPost  error
	}{
		{
			Name:         "success",
			PostId:       1,
			UserId:       1,
			ExpectedRes:  p,
			ExpectedPost: p,
		},
		{
			Name:        "failed post not found",
			PostId:      1,
			UserId:      1,
			ExpectedErr: post.ErrorPostNotFound,
			ErrRepoPost: post.ErrorPostNotFound,
		},
		{
			Name:   "failed user unauthorized",
			PostId: 1,
			UserId: 1,
			ExpectedPost: post.Post{
				Id:      1,
				UserId:  4,
				Title:   "test",
				Content: "test",
			},
			ExpectedErr: post.ErrorUnauthorized,
		},
		{
			Name:        "failed",
			PostId:      1,
			UserId:      1,
			ExpectedErr: post.ErrDatabaseFailure,
			ErrRepoPost: post.ErrDatabaseFailure,
		},
	}

	for _, tc := range testCases {
		r := mocks.NewPostRepository(t)
		i := post.NewInteractor(r)

		r.On("FindByPostId", tc.PostId).Return(tc.ExpectedPost, tc.ErrRepoPost)
		res, err := i.GetPost(tc.UserId, tc.PostId)
		assert.Equal(t, tc.ExpectedErr, err)
		assert.Equal(t, tc.ExpectedRes, res)
	}
}

func TestGetPostByUserId(t *testing.T) {
	testCases := []struct {
		Name        string
		UserId      int
		ExpectedRes []post.Post
		ExpectedErr error
	}{
		{
			Name:   "success",
			UserId: 1,
			ExpectedRes: []post.Post{
				{
					Id:      1,
					UserId:  1,
					Title:   "test",
					Content: "test",
				},
				{
					Id:      2,
					UserId:  1,
					Title:   "test",
					Content: "test",
				},
			},
		},
		{
			Name:        "failed post not found",
			UserId:      8,
			ExpectedErr: post.ErrorPostNotFound,
		},
		{
			Name:        "failed",
			UserId:      2,
			ExpectedErr: post.ErrDatabaseFailure,
		},
	}

	for _, tc := range testCases {
		r := mocks.NewPostRepository(t)
		i := post.NewInteractor(r)

		r.On("FindByUserId", tc.UserId).Return(tc.ExpectedRes, tc.ExpectedErr)
		res, err := i.GetPostByUserId(tc.UserId)
		assert.Equal(t, tc.ExpectedErr, err)
		assert.Equal(t, tc.ExpectedRes, res)
	}
}

func TestUpdatePost(t *testing.T) {
	p := post.Post{
		Id:      1,
		UserId:  1,
		Title:   "test update",
		Content: "test update",
	}

	testCases := []struct {
		Name         string
		PostId       int
		userId       int
		Post         post.Post
		ExpectedRes  post.Post
		ExpectedErr  error
		CallRepoPost bool
	}{
		{
			Name:         "success",
			userId:       1,
			PostId:       1,
			Post:         p,
			ExpectedRes:  p,
			CallRepoPost: true,
		},
		{
			Name:         "failed user unauthorized",
			userId:       9,
			PostId:       1,
			Post:         p,
			ExpectedErr:  post.ErrorUnauthorized,
			CallRepoPost: false,
		},
		{
			Name:         "failed",
			userId:       1,
			PostId:       1,
			Post:         p,
			CallRepoPost: true,
			ExpectedErr:  post.ErrDatabaseFailure,
		},
	}

	for _, tc := range testCases {
		r := mocks.NewPostRepository(t)
		i := post.NewInteractor(r)

		if tc.CallRepoPost {
			r.On("Update", tc.PostId, tc.Post).Return(tc.ExpectedErr)
		}

		res, err := i.UpdatePost(tc.PostId, tc.userId, tc.Post)
		assert.Equal(t, tc.ExpectedErr, err)
		assert.Equal(t, tc.ExpectedRes, res)
	}
}

func TestDeletePost(t *testing.T) {
	p := post.Post{
		Id:      1,
		UserId:  1,
		Title:   "test",
		Content: "test",
	}

	testCases := []struct {
		Name        string
		PostID      int
		UserID      int
		Post        post.Post
		ErrPostRepo error
		ExpectedErr error
		IsValidUser bool
	}{
		{
			Name:        "success",
			PostID:      1,
			UserID:      1,
			Post:        p,
			IsValidUser: true,
		},
		{
			Name:        "failed post not found",
			PostID:      2,
			UserID:      1,
			Post:        post.Post{},
			ErrPostRepo: post.ErrorPostNotFound,
			ExpectedErr: post.ErrorPostNotFound,
		},
		{
			Name:        "failed user unauthorized",
			PostID:      1,
			UserID:      8,
			Post:        p,
			ExpectedErr: post.ErrorUnauthorized,
		},
		{
			Name:        "failed",
			PostID:      1,
			UserID:      1,
			Post:        p,
			ExpectedErr: post.ErrDatabaseFailure,
			IsValidUser: true,
		},
	}

	for _, tc := range testCases {
		r := mocks.NewPostRepository(t)
		i := post.NewInteractor(r)

		r.On("FindByPostId", tc.PostID).Return(tc.Post, tc.ErrPostRepo)
		if tc.IsValidUser {
			r.On("Delete", tc.PostID).Return(tc.ExpectedErr)
		}

		err := i.DeletePost(tc.PostID, tc.UserID)
		assert.Equal(t, tc.ExpectedErr, err)
	}
}
