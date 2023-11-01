package handler

import (
	"go-post/internal/helper"
	"go-post/internal/post"
	"go-post/internal/response"
	"go-post/internal/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type postHandler struct {
	userInteractor user.Interactor
	postInteractor post.Interactor
}

func NewPostHandler(
	userInteractor user.Interactor,
	postInteractor post.Interactor,
) postHandler {
	return postHandler{
		userInteractor: userInteractor,
		postInteractor: postInteractor,
	}
}

func (h *postHandler) CreatePost(c *gin.Context) {
	var input post.InputCreatePost

	if err := c.BindJSON(&input); err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", err.Error(), c, false)
		return
	}

	userId := c.MustGet("user").(int)

	post := post.Post{
		UserId:  userId,
		Title:   input.Title,
		Content: input.Content,
	}

	if httpCode, err := h.postInteractor.CreatePost(post); err != nil {
		helper.GenerateResponseAPI(httpCode, "error", err.Error(), c, false)
		return
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", "success", c, false)
}

func (h *postHandler) GetPost(c *gin.Context) {
	id := c.Param("id")
	userId := c.MustGet("user").(int)

	postId, err := strconv.Atoi(id)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error in convert id", err.Error(), c, false)
		return
	}

	post, httpCode, err := h.postInteractor.GetPost(userId, postId)
	if err != nil {
		helper.GenerateResponseAPI(httpCode, "error", err.Error(), c, false)
		return
	}

	user, err := h.userInteractor.GetUserById(userId)
	if err != nil {
		helper.GenerateResponseAPI(httpCode, "error", err.Error(), c, false)
		return
	}

	postResponse := response.PostResponse{
		Post: post,
		User: user,
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", postResponse, c, false)
}

func (h *postHandler) DeletePost(c *gin.Context) {
	userId := c.MustGet("user").(int)
	id := c.Param("id")

	postId, err := strconv.Atoi(id)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error in convert id", err.Error(), c, false)
		return
	}

	httpCode, err := h.postInteractor.DeletePost(postId, userId)
	if err != nil {
		helper.GenerateResponseAPI(httpCode, "error", err.Error(), c, false)
		return
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", "success", c, false)
}

func (h *postHandler) UpdatePost(c *gin.Context) {
	userId := c.MustGet("user").(int)
	id := c.Param("id")

	postId, err := strconv.Atoi(id)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error in convert id", err.Error(), c, false)
		return
	}

	var input post.InputUpdatePost
	if err := c.BindJSON(&input); err != nil {
		errsBinding := helper.ErrorBindingFormatter(err)
		helper.GenerateResponseAPI(http.StatusBadRequest, "error binding", errsBinding, c, false)
		return
	}

	post, httpCode, err := h.postInteractor.GetPost(userId, postId)
	if err != nil {
		helper.GenerateResponseAPI(httpCode, "error", err.Error(), c, false)
		return
	}

	post.Content = input.Content
	post.Title = input.Title

	post, httpCode, err = h.postInteractor.UpdatePost(postId, userId, post)
	if err != nil {
		helper.GenerateResponseAPI(httpCode, "error", err.Error(), c, false)
		return
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", post, c, false)
}
