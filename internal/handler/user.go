package handler

import (
	"go-post/internal/helper"
	"go-post/internal/post"
	"go-post/internal/response"
	"go-post/internal/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type userHandler struct {
	userRepo user.UserRepository
	postRepo post.PostRepository
}

func NewUserHandler(userRepo user.UserRepository, postRepo post.PostRepository) *userHandler {
	return &userHandler{
		userRepo: userRepo,
		postRepo: postRepo,
	}
}

func (h *userHandler) SignUp(c *gin.Context) {
	var input user.InputUserSignUp

	if err := c.BindJSON(&input); err != nil {
		errsBinding := helper.ErrorBindingFormatter(err)
		helper.GenerateResponseAPI(http.StatusBadRequest, "error binding", errsBinding, c, false)
		return
	}

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", err, c, false)
		return
	}

	var newUser user.User
	newUser.Email = input.Email
	newUser.Password = string(passwordHashed)

	err = h.userRepo.Save(newUser)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", err, c, false)
		return
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", "success", c, false)
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.InputUserLogin

	if err := c.BindJSON(&input); err != nil {
		errsBinding := helper.ErrorBindingFormatter(err)
		helper.GenerateResponseAPI(http.StatusBadRequest, "error binding", errsBinding, c, false)
		return
	}

	user, err := h.userRepo.FindByEmail(input.Email)
	if err != nil && user.Id == 0 {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", "email not registered", c, false)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", "wrong password", c, false)
		return
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", user, c, false)
}

func (h *userHandler) GetUserWithPosts(c *gin.Context) {
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

	userWithPosts := response.UserWithPostsResponse{
		User: user,
		Post: posts,
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", userWithPosts, c, false)
}
