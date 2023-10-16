package user_test

import (
	"go-post/internal/database"
	"go-post/internal/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveUser(t *testing.T) {
	db, err := database.NewConnection("../../.env")
	assert.NoError(t, err)
	r := user.NewRepositoryUser(db)

	u := user.User{
		Email:    "test@gmail.com",
		Password: "1234567",
	}

	res, err := r.Save(u)
	assert.NoError(t, err)
	assert.Equal(t, u.Email, res.Email)
	assert.Equal(t, u.Password, res.Password)

	defer func() {
		db.Exec("TRUNCATE users")
	}()
}

func TestFindByEmail(t *testing.T) {
	db, err := database.NewConnection("../../.env")
	assert.NoError(t, err)
	r := user.NewRepositoryUser(db)

	u := user.User{
		Email:    "test@gmail.com",
		Password: "1234567",
	}

	_, err = r.Save(u)
	assert.NoError(t, err)

	user, err := r.FindByEmail(u.Email)
	assert.NoError(t, err)
	assert.Equal(t, u.Email, user.Email)
	assert.Equal(t, u.Password, user.Password)

	defer func() {
		db.Exec("TRUNCATE users")
	}()
}

func TestFindById(t *testing.T) {
	db, err := database.NewConnection("../../.env")
	assert.NoError(t, err)
	r := user.NewRepositoryUser(db)

	u := user.User{
		Email:    "test@gmail.com",
		Password: "1234567",
	}

	res, err := r.Save(u)
	assert.NoError(t, err)

	user, err := r.FindById(res.Id)
	assert.NoError(t, err)
	assert.Equal(t, res.Id, user.Id)
	assert.Equal(t, res.Email, user.Email)
	assert.Equal(t, res.Password, user.Password)

	defer func() {
		db.Exec("TRUNCATE users")
	}()
}
