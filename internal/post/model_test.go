package post_test

import (
	"go-post/internal/post"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModelPostStruct(t *testing.T) {
	p := post.Post{
		Id:      1,
		UserId:  3,
		Title:   "test",
		Content: "test",
	}

	assert.Equal(t, 1, p.Id)
	assert.Equal(t, 3, p.UserId)
	assert.Equal(t, "test", p.Title)
	assert.Equal(t, "test", p.Content)
}

func TestInputCreatePostStruct(t *testing.T) {
	i := post.InputCreatePost{
		Title:   "test",
		Content: "test",
	}

	assert.Equal(t, "test", i.Title)
	assert.Equal(t, "test", i.Content)
}

func TestInputUpdatePostStruct(t *testing.T) {
	i := post.InputUpdatePost{
		UserId:  3,
		Title:   "test",
		Content: "test",
	}

	assert.Equal(t, 3, i.UserId)
	assert.Equal(t, "test", i.Title)
	assert.Equal(t, "test", i.Content)
}
