package database_test

import (
	"go-post/internal/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnection(t *testing.T) {
	db, err := database.NewConnection("../../.env")
	assert.NoError(t, err)
	assert.NotNil(t, db)
}
