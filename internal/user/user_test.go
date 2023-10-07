package user_test

import (
	"go-post/internal/user"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestUserStruct(t *testing.T) {
	user := user.User{
		Id: 1,
		Email: "test@gmail.com",
		Password: "jdbjhfbdjhd-shdb-12",
	}

	require.Equal(t, 1, user.Id)
	require.Equal(t, "test@gmail.com", user.Email)
	require.Equal(t, "jdbjhfbdjhd-shdb-12", user.Password)
}
