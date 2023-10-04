package handler

import (
	"go-post/internal/helper"
	"go-post/internal/model"
	"go-post/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type postHandler struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

func NewPostHandler(postRepo repository.PostRepository, userRepo repository.UserRepository) postHandler {
	return postHandler{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

func (h *postHandler) CreatePost(c *gin.Context) {
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

func (h *postHandler) GetPostDetail(c *gin.Context) {
	id := c.Param("id")

	postId, err := strconv.Atoi(id)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error in convert id", err.Error(), c, false)
		return
	}

	post, err := h.postRepo.FindByPostId(postId)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusInternalServerError, "error", err.Error(), c, false)
		return
	}

	user, err := h.userRepo.FindById(post.UserId)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusInternalServerError, "error", err.Error(), c, false)
		return
	}

	postDetail := model.PostDetailReponse{
		Post: post,
		User: user,
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", postDetail, c, false)
}

func (h *postHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")

	postId, err := strconv.Atoi(id)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error in convert id", err.Error(), c, false)
		return
	}

	post, err := h.postRepo.FindByPostId(postId)
	if err != nil && post.Id == 0 {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", err.Error(), c, false)
		return
	}

	if err := h.postRepo.Delete(postId); err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", err.Error(), c, false)
		return
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", "success", c, false)
}
