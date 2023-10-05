package response

import (
	"go-post/internal/post"
	"go-post/internal/user"
)

type PostWithUserReponse struct {
	Post post.Post `json:"post"`
	User user.User `json:"user"`
}

type UserWithPostsResponse struct {
	User user.User   `json:"user"`
	Post []post.Post `json:"posts"`
}
