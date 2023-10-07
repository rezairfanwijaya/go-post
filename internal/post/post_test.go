package post_test

import (
	"go-post/internal/post"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostStruct(t *testing.T) {
	post := post.Post{
		Id:      1,
		UserId:  2,
		Title:   "test",
		Content: "detail test content",
	}

	assert.Equal(t, 1, post.Id)
	assert.Equal(t, 2, post.UserId)
	assert.Equal(t, "test", post.Title)
	assert.Equal(t, "detail test content", post.Content)
}
