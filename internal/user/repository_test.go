package user_test

import (
	"database/sql"
	"errors"
	"go-post/internal/user"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func NewConnection() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return db, mock, err
	}

	return db, mock, nil
}

func TestSaveUser(t *testing.T) {
	db, mock, err := NewConnection()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := user.NewRepositoryUser(db)

	testCases := []struct {
		Name      string
		Param     user.User
		WantError bool
	}{
		{
			Name: "success",
			Param: user.User{
				Email:    "andi@gmail.com",
				Password: "jdbfjhb4hsbjd",
			},
			WantError: false,
		},
		{
			Name:      "failed",
			Param:     user.User{},
			WantError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if !testCase.WantError {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (email, password) VALUES ($1, $2)")).WithArgs(testCase.Param.Email, testCase.Param.Password).WillReturnResult(sqlmock.NewResult(1, 1))

				err := repo.Save(testCase.Param)
				assert.Nil(t, err)
			} else {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (email, password) VALUES ($1, $2)")).WithArgs(testCase.Param.Email, testCase.Param.Password).WillReturnError(errors.New("failed"))

				err := repo.Save(testCase.Param)
				assert.NotNil(t, err)
			}
		})
	}
}