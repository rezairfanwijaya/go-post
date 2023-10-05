package article

import (
	"go-post/internal/helper"
	"go-post/internal/post"
	"go-post/internal/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type articleHandler struct {
	userRepo user.UserRepository
	postRepo post.PostRepository
}

func NewArticleHandler(userRepo user.UserRepository, postRepo post.PostRepository) articleHandler {
	return articleHandler{
		userRepo: userRepo,
		postRepo: postRepo,
	}
}

func (h *articleHandler) GetArticlesDetail(c *gin.Context) {
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

	articleDetail := ArticleDetail{
		Post: post,
		User: user,
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", articleDetail, c, false)
}

func (h *articleHandler) GetArticlesByUser(c *gin.Context) {
	id := c.Param("id")

	userId, err := strconv.Atoi(id)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error in convert id", err.Error(), c, false)
		return
	}

	user, err := h.userRepo.FindById(userId)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", err.Error(), c, false)
		return
	}

	posts, err := h.postRepo.FindByUserId(user.Id)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", err.Error(), c, false)
		return
	}

	articlesWithUser := ArticlesWithUser{
		User:  user,
		Posts: posts,
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", articlesWithUser, c, false)
}
