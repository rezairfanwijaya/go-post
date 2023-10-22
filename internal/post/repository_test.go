package post_test

import (
	"errors"
	"go-post/internal/database"
	"go-post/internal/post"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	db, err := database.NewConnection("../../.env")
	assert.NoError(t, err)
	r := post.NewPostRepository(db)

	p := post.Post{
		UserId:  1,
		Title:   "test",
		Content: "test",
	}

	np, err := r.Save(p)
	assert.NoError(t, err)
	assert.Equal(t, p.UserId, np.UserId)
	assert.Equal(t, p.Title, np.Title)
	assert.Equal(t, p.Content, np.Content)

	defer func() {
		db.Exec("DELETE FROM posts WHERE id = $1", np.Id)
	}()
}

func TestFindByPostId(t *testing.T) {
	db, err := database.NewConnection("../../.env")
	assert.NoError(t, err)
	r := post.NewPostRepository(db)

	testCases := []struct {
		Name        string
		Post        post.Post
		ErrExpected error
	}{
		{
			Name: "success",
			Post: post.Post{
				UserId:  1,
				Title:   "test",
				Content: "test",
			},
		},
		{
			Name: "failed_not_found",
			Post: post.Post{},
		},
	}

	for _, tc := range testCases {
		np, err := r.Save(tc.Post)
		assert.NoError(t, err)

		p, err := r.FindByPostId(np.Id)
		assert.NoError(t, err)

		assert.Equal(t, tc.Post.Title, p.Title)
		assert.Equal(t, tc.Post.Content, p.Content)

		defer func() {
			db.Exec("DELETE FROM posts WHERE id = $1", np.Id)
		}()
	}
}

func TestFindByUserId(t *testing.T) {
	db, err := database.NewConnection("../../.env")
	assert.NoError(t, err)

	p := post.Post{
		UserId:  3,
		Title:   "this is title",
		Content: "this content",
	}

	testCases := []struct {
		Name         string
		Post         post.Post
		PostResponse []post.Post
	}{
		{
			Name: "success",
			Post: p,
			PostResponse: []post.Post{
				p,
			},
		},
		{
			Name:         "not_found",
			Post:         post.Post{},
			PostResponse: []post.Post{},
		},
	}

	for _, tc := range testCases {
		r := post.NewPostRepository(db)
		newPost, err := r.Save(tc.Post)
		assert.NoError(t, err)
		assert.Equal(t, tc.Post.Title, newPost.Title)
		func() {
			db.Exec("DELETE FROM posts WHERE id=$1", newPost.Id)
		}()
	}

}

func TestDelete(t *testing.T) {
	db, err := database.NewConnection("../../.env")
	assert.NoError(t, err)

	p := post.Post{
		UserId:  3,
		Title:   "this is title",
		Content: "this content",
	}

	r := post.NewPostRepository(db)
	newPost, err := r.Save(p)
	assert.NoError(t, err)
	err = r.Delete(newPost.Id)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := post.NewPostRepository(db)

	testCases := []struct {
		Name    string
		PostId  int
		Post    post.Post
		WantErr bool
	}{
		{
			Name:   "success",
			PostId: 1,
			Post: post.Post{
				Id:      1,
				Title:   "title update",
				Content: "content update",
			},
			WantErr: false,
		},
		{
			Name:    "false",
			PostId:  0,
			Post:    post.Post{},
			WantErr: true,
		},
	}

	for _, testCase := range testCases {
		if !testCase.WantErr {
			mock.ExpectExec(regexp.QuoteMeta("UPDATE posts SET title = $1, content = $2 WHERE id = $3")).WithArgs(testCase.Post.Title, testCase.Post.Content, testCase.Post.Id).WillReturnResult(sqlmock.NewResult(0, 1))

			err := repo.Update(testCase.PostId, testCase.Post)
			assert.Nil(t, err)
		} else {
			mock.ExpectExec(regexp.QuoteMeta("UPDATE posts SET title = $1, content = $2 WHERE id = $3")).WithArgs(testCase.Post.Title, testCase.Post.Content, testCase.Post.Id).WillReturnError(errors.New("failed"))

			err := repo.Update(testCase.PostId, testCase.Post)
			assert.NotNil(t, err)
		}
	}
}
