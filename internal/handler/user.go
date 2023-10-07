package handler

import (
	"go-post/internal/auth"
	"go-post/internal/helper"
	"go-post/internal/post"
	"go-post/internal/response"
	"go-post/internal/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type userHandler struct {
	postInteractor post.Interactor
	userInteractor user.Interactor
	authService    auth.AuthService
}

func NewUserHandler(
	postInteractor post.Interactor,
	userInteractor user.Interactor,
	authService auth.AuthService,
) *userHandler {
	return &userHandler{
		postInteractor: postInteractor,
		userInteractor: userInteractor,
		authService:    authService,
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

	httpCode, err := h.userInteractor.CreateUser(newUser)
	if err != nil {
		helper.GenerateResponseAPI(httpCode, "error", err, c, false)
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

	user, httpCode, err := h.userInteractor.GetUserByEmail(input.Email)
	if err != nil {
		helper.GenerateResponseAPI(httpCode, "error", err.Error(), c, false)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", "wrong password", c, false)
		return
	}

	token, err := h.authService.GenerateToken(user.Id)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusInternalServerError, "error generate token", err.Error(), c, false)
		return
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", token, c, false)
}

func (h *userHandler) GetUserWithPosts(c *gin.Context) {
	userId := c.MustGet("user").(int)

	posts, httpCode, err := h.postInteractor.GetPostByUserId(userId)
	if err != nil {
		helper.GenerateResponseAPI(httpCode, "error", err.Error(), c, false)
		return
	}

	user, httpCode, err := h.userInteractor.GetUserById(userId)
	if err != nil {
		helper.GenerateResponseAPI(httpCode, "error", err.Error(), c, false)
		return
	}

	userWithPosts := response.UserWithPostsResponse{
		User: user,
		Post: posts,
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", userWithPosts, c, false)
}
