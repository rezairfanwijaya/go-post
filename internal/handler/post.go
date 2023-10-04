package handler

import "go-post/internal/repository"

type postHandler struct {
	postRepo repository.PostRepository
}

func NewPostHandler(postRepo repository.PostRepository) postHandler {
	return postHandler{
		postRepo: postRepo,
	}
}
