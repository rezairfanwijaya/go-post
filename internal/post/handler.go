package post

import (
	"go-post/internal/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type postHandler struct {
	postRepo PostRepository
}

func NewPostHandler(postRepo PostRepository) postHandler {
	return postHandler{
		postRepo: postRepo,
	}
}

func (h *postHandler) CreatePost(c *gin.Context) {
	var input InputCreatePost

	if err := c.BindJSON(&input); err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", err.Error(), c, false)
		return
	}

	post := Post{
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

func (h *postHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")

	postId, err := strconv.Atoi(id)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error in convert id", err.Error(), c, false)
		return
	}

	var input InputUpdatePost
	if err := c.BindJSON(&input); err != nil {
		errsBinding := helper.ErrorBindingFormatter(err)
		helper.GenerateResponseAPI(http.StatusBadRequest, "error binding", errsBinding, c, false)
		return
	}

	post, err := h.postRepo.FindByPostId(postId)
	if err != nil && post.Id == 0 {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", err.Error(), c, false)
		return
	}

	post.Content = input.Content
	post.Title = input.Title

	err = h.postRepo.Update(postId, post)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusInternalServerError, "error", err.Error(), c, false)
		return
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", post, c, false)
}
