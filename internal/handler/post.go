package handler

import (
	"errors"
	"go-post/internal/helper"
	"go-post/internal/post"
	"go-post/internal/response"
	"go-post/internal/user"
	"log"
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
		log.Printf("failed to binding input, err: %s", err)
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", "Invalid input", c, false)
		return
	}

	userId := c.MustGet("user").(int)

	post := post.Post{
		UserId:  userId,
		Title:   input.Title,
		Content: input.Content,
	}

	p, err := h.postInteractor.CreatePost(post)
	if err != nil {
		log.Printf("failed to create new post, err: %s", err)
		helper.GenerateResponseAPI(http.StatusInternalServerError, "error", "Failed create post", c, false)
		return
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", p, c, false)
}

func (h *postHandler) GetPost(c *gin.Context) {
	id := c.Param("id")
	userId := c.MustGet("user").(int)

	postId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("failed convert to integer, err: %s", err)
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", "Invalid input", c, false)
		return
	}

	p, err := h.postInteractor.GetPost(userId, postId)
	if err != nil {
		if errors.Is(err, post.ErrorPostNotFound) {
			log.Printf("failed to get post, postID: %d, err: %s", postId, err)
			helper.GenerateResponseAPI(http.StatusNotFound, "error", "Post not found", c, false)
			return
		}

		if errors.Is(err, post.ErrorUnauthorized) {
			log.Printf("failed to get post, postID: %d, userID: %d, err: %s", postId, userId, err)
			helper.GenerateResponseAPI(http.StatusUnauthorized, "error", "Unauthorized", c, false)
			return
		}

		log.Printf("failed to get post, postID: %d, err: %s", postId, err)
		helper.GenerateResponseAPI(http.StatusNotFound, "error", "Unknown error occurred", c, false)
		return
	}

	u, err := h.userInteractor.GetUserById(userId)
	if err != nil {
		if errors.Is(err, user.ErrorUserNotFound) {
			helper.GenerateResponseAPI(http.StatusNotFound, "error", "User not found", c, false)
			return
		}

		helper.GenerateResponseAPI(http.StatusInternalServerError, "error", "Unknown error occurred", c, false)
		return
	}

	postResponse := response.PostResponse{
		Post: p,
		User: u,
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", postResponse, c, false)
}

func (h *postHandler) DeletePost(c *gin.Context) {
	userId := c.MustGet("user").(int)
	id := c.Param("id")

	postId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("failed convert to integer, err: %s", err)
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", "Invalid input", c, false)
		return
	}

	err = h.postInteractor.DeletePost(postId, userId)
	if err != nil {
		if errors.Is(err, post.ErrorPostNotFound) {
			log.Printf("failed to delete post, postID: %d, err: %s", postId, err)
			helper.GenerateResponseAPI(http.StatusNotFound, "error", "Post not found", c, false)
			return
		}

		if errors.Is(err, post.ErrorUnauthorized) {
			log.Printf("failed to delete post, postID: %d, userID: %d, err: %s", postId, userId, err)
			helper.GenerateResponseAPI(http.StatusUnauthorized, "error", "Unauthorized", c, false)
			return
		}

		helper.GenerateResponseAPI(http.StatusInternalServerError, "error", "Unknown error occurred", c, false)
		return
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", "success", c, false)
}

func (h *postHandler) UpdatePost(c *gin.Context) {
	userId := c.MustGet("user").(int)
	id := c.Param("id")

	postId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("failed convert to integer, err: %s", err)
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", "Invalid input", c, false)
		return
	}

	var input post.InputUpdatePost
	if err := c.BindJSON(&input); err != nil {
		log.Printf("failed to binding input, err: %s", err)
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", "Invalid input", c, false)
		return
	}

	p, err := h.postInteractor.GetPost(userId, postId)
	if err != nil {
		if errors.Is(err, post.ErrorPostNotFound) {
			log.Printf("failed to get post, postID: %d, err: %s", postId, err)
			helper.GenerateResponseAPI(http.StatusNotFound, "error", "Post not found", c, false)
			return
		}

		if errors.Is(err, post.ErrorUnauthorized) {
			log.Printf("failed to update post, postID: %d, userID: %d, err: %s", postId, userId, err)
			helper.GenerateResponseAPI(http.StatusUnauthorized, "error", "Unauthorized", c, false)
			return
		}

		log.Printf("failed to get post, postID: %d, err: %s", postId, err)
		helper.GenerateResponseAPI(http.StatusNotFound, "error", "Unknown error occurred", c, false)
		return
	}

	p.Content = input.Content
	p.Title = input.Title

	res, err := h.postInteractor.UpdatePost(postId, userId, p)
	if err != nil {
		if errors.Is(err, post.ErrorUnauthorized) {
			log.Printf("failed to update post, postID: %d, userID: %d, err: %s", postId, userId, err)
			helper.GenerateResponseAPI(http.StatusUnauthorized, "error", "Unauthorized", c, false)
			return
		}

		helper.GenerateResponseAPI(http.StatusInternalServerError, "error", "Unknow error occurred", c, false)
		return
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", res, c, false)
}
