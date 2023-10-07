package post_test

import (
	"errors"
	"go-post/internal/post"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := post.NewPostRepository(db)

	testCases := []struct {
		Name    string
		Param   post.Post
		WantErr bool
	}{
		{
			Name: "success",
			Param: post.Post{
				Id:      1,
				UserId:  2,
				Title:   "test",
				Content: "detail content",
			},
			WantErr: false,
		},
		{
			Name:    "failed",
			Param:   post.Post{},
			WantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if testCase.WantErr {
				mock.ExpectExec(regexp.QuoteMeta("INSERT into posts (user_id, title, content) VALUES ($1, $2, $3)")).
					WithArgs(testCase.Param.UserId, testCase.Param.Title, testCase.Param.Content).
					WillReturnError(errors.New("failed"))

				err := repo.Save(testCase.Param)
				assert.NotNil(t, err)
			} else {
				mock.ExpectExec(regexp.QuoteMeta("INSERT into posts (user_id, title, content) VALUES ($1, $2, $3)")).
					WithArgs(testCase.Param.UserId, testCase.Param.Title, testCase.Param.Content).WillReturnResult(sqlmock.NewResult(1, 1))

				err := repo.Save(testCase.Param)
				assert.Nil(t, err)
			}
		})
	}
}

func TestFindByPostId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := post.NewPostRepository(db)

	testCases := []struct {
		Name        string
		Id          int
		Expectation post.Post
		WantErr     bool
	}{
		{
			Name: "success",
			Id:   1,
			Expectation: post.Post{
				Id:      1,
				UserId:  2,
				Title:   "test",
				Content: "test detail in content",
			},
			WantErr: false,
		},
		{
			Name:        "not found",
			Id:          4,
			Expectation: post.Post{},
			WantErr:     false,
		},
		{
			Name:        "failed",
			Id:          0,
			Expectation: post.Post{},
			WantErr:     true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if !testCase.WantErr {
				rows := mock.NewRows([]string{"id", "user_id", "title", "content"}).
					AddRow(testCase.Expectation.Id, testCase.Expectation.UserId, testCase.Expectation.Title, testCase.Expectation.Content)

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM posts WHERE id = $1")).WithArgs(testCase.Id).
					WillReturnRows(rows)

				post, _ := repo.FindByPostId(testCase.Id)
				assert.Equal(t, testCase.Expectation, post)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM posts WHERE id = $1")).WithArgs(testCase.Id).
					WillReturnError(errors.New("failed"))

				_, err := repo.FindByPostId(testCase.Id)
				assert.NotNil(t, err)
			}
		})
	}
}

func TestFindByUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := post.NewPostRepository(db)

	testCases := []struct {
		Name        string
		Id          int
		Expectation []post.Post
		WantErr     bool
	}{
		{
			Name: "success",
			Id:   2,
			Expectation: []post.Post{
				{
					Id:      1,
					Title:   "test",
					Content: "test detail in content",
				}, {
					Id:      2,
					Title:   "test content sport",
					Content: "test detail in content sport",
				},
			},
			WantErr: false,
		},
		{
			Name:        "failed",
			Id:          0,
			Expectation: []post.Post{},
			WantErr:     true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if !testCase.WantErr {
				rows := mock.NewRows([]string{"id", "title", "content"}).
					AddRow(testCase.Expectation[0].Id, testCase.Expectation[0].Title, testCase.Expectation[0].Content).
					AddRow(testCase.Expectation[1].Id, testCase.Expectation[1].Title, testCase.Expectation[1].Content)

				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, content FROM posts WHERE user_id = $1")).WithArgs(testCase.Id).
					WillReturnRows(rows)

				posts, _ := repo.FindByUserId(testCase.Id)
				assert.Equal(t, testCase.Expectation, posts)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM posts WHERE user_id = $1")).WithArgs(testCase.Id).
					WillReturnError(errors.New("failed"))

				_, err := repo.FindByUserId(testCase.Id)
				assert.NotNil(t, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := post.NewPostRepository(db)

	testCases := []struct {
		Name    string
		Param   int
		WantErr bool
	}{
		{
			Name:    "success",
			Param:   1,
			WantErr: false,
		},
		{
			Name:    "failed",
			Param:   0,
			WantErr: true,
		},
	}

	for _, testCase := range testCases {
		if !testCase.WantErr {
			mock.ExpectExec(regexp.QuoteMeta("DELETE FROM posts WHERE id = $1")).WithArgs(testCase.Param).
				WillReturnResult(sqlmock.NewResult(0, 1))

			err := repo.Delete(testCase.Param)
			assert.Nil(t, err)
		} else {
			mock.ExpectExec(regexp.QuoteMeta("DELETE FROM posts WHERE id = $1")).WithArgs(testCase.Param).
				WillReturnError(errors.New("failed"))

			err := repo.Delete(testCase.Param)
			assert.NotNil(t, err)
		}
	}
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
