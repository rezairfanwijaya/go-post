package article

import (
	"go-post/internal/post"
	"go-post/internal/user"
)

type ArticleDetail struct {
	Post post.Post `json:"post"`
	User user.User `json:"user"`
}

type ArticlesWithUser struct {
	User  user.User   `json:"user"`
	Posts []post.Post `json:"posts"`
}
