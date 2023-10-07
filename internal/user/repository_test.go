package user_test

import (
	"errors"
	"go-post/internal/user"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSaveUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := user.NewRepositoryUser(db)

	testCases := []struct {
		Name    string
		Param   user.User
		WantErr bool
	}{
		{
			Name: "success",
			Param: user.User{
				Id:       1,
				Email:    "test@gmail.com",
				Password: "jdjfbjdfbdj",
			},
			WantErr: false,
		},
		{
			Name:    "failed",
			Param:   user.User{},
			WantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if testCase.WantErr {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (email, password) VALUES ($1, $2)")).
					WithArgs(testCase.Param.Email, testCase.Param.Password).WillReturnError(errors.New("failed"))

				err := repo.Save(testCase.Param)
				assert.NotNil(t, err)
			} else {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (email, password) VALUES ($1, $2)")).
					WithArgs(testCase.Param.Email, testCase.Param.Password).WillReturnResult(sqlmock.NewResult(1, 1))

				err := repo.Save(testCase.Param)
				assert.Nil(t, err)
			}
		})
	}

}

func TestFindByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := user.NewRepositoryUser(db)

	testCases := []struct {
		Name        string
		Email       string
		Expectation user.User
		WantErr     bool
	}{
		{
			Name:  "success",
			Email: "test@gmail.com",
			Expectation: user.User{
				Id:       1,
				Email:    "test@gmail.com",
				Password: "dfjbgdjhfbg54",
			},
			WantErr: false,
		},
		{
			Name:        "not found",
			Email:       "agnes@gmail.com",
			Expectation: user.User{},
			WantErr:     false,
		},
		{
			Name:        "failed",
			Email:       "",
			Expectation: user.User{},
			WantErr:     true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if !testCase.WantErr {
				rows := mock.NewRows([]string{"id", "email", "password"}).
					AddRow(testCase.Expectation.Id, testCase.Expectation.Email, testCase.Expectation.Password)

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE email = $1")).WithArgs(testCase.Email).
					WillReturnRows(rows)

				user, _ := repo.FindByEmail(testCase.Email)
				assert.Equal(t, testCase.Expectation, user)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE email = $1")).WithArgs(testCase.Email).
					WillReturnError(errors.New("failed"))

				_, err := repo.FindByEmail(testCase.Email)
				assert.NotNil(t, err)
			}
		})
	}
}

func TestFindById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := user.NewRepositoryUser(db)

	testCases := []struct {
		Name        string
		Id          int
		Expectation user.User
		WantErr     bool
	}{
		{
			Name: "success",
			Id:   1,
			Expectation: user.User{
				Id:       1,
				Email:    "test@gmail.com",
				Password: "dfjbgdjhfbg54",
			},
			WantErr: false,
		},
		{
			Name:        "not found",
			Id:          89,
			Expectation: user.User{},
			WantErr:     false,
		},
		{
			Name:        "failed",
			Id:          0,
			Expectation: user.User{},
			WantErr:     true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if !testCase.WantErr {
				rows := mock.NewRows([]string{"id", "email", "password"}).
					AddRow(testCase.Expectation.Id, testCase.Expectation.Email, testCase.Expectation.Password)

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE id = $1")).WithArgs(testCase.Id).
					WillReturnRows(rows)

				user, _ := repo.FindById(testCase.Id)
				assert.Equal(t, testCase.Expectation, user)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE id = $1")).WithArgs(testCase.Id).
					WillReturnError(errors.New("failed"))

				_, err := repo.FindById(testCase.Id)
				assert.NotNil(t, err)
			}
		})
	}
}
