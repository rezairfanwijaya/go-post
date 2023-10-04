package handler

import (
	"go-post/internal/helper"
	"go-post/internal/model"
	"go-post/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type postHandler struct {
	postRepo repository.PostRepository
}

func NewPostHandler(postRepo repository.PostRepository) postHandler {
	return postHandler{
		postRepo: postRepo,
	}
}

func (h *postHandler) NewPost(c *gin.Context) {
	var input model.InputCreatePost

	if err := c.BindJSON(&input); err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", err.Error(), c, false)
		return
	}

	post := model.Post{
		UserId:  input.UserId,
		Title:   input.Title,
		Content: input.Content,
	}

	if err := h.postRepo.Save(post); err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", err.Error(), c, false)
		return
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", "success", c, false)
}
