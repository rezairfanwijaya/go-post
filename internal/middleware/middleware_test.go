package middleware_test

import (
	"go-post/internal/auth"
	"go-post/internal/middleware"
	"go-post/internal/user"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	auth := auth.NewAuthService()
	userRepo := user.NewRepositoryUser(db)
	userInteractor := user.NewInteractor(userRepo)

	handler := middleware.AuthFunc(auth, userInteractor)
	assert.NotNil(t, handler)
}
