package handler

import (
	"go-post/internal/auth"
	"go-post/internal/helper"
	"go-post/internal/post"
	"go-post/internal/response"
	"go-post/internal/user"
	"log"
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
		log.Printf("failed to binding json input, err: %v", errsBinding)
		helper.GenerateResponseAPI(http.StatusBadRequest, "error binding", "Invalid input", c, false)
		return
	}

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		log.Printf("failed to generate hash password, err: %v", err)
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", "Invalid password", c, false)
		return
	}

	var newUser user.User
	newUser.Email = input.Email
	newUser.Password = string(passwordHashed)

	user, err := h.userInteractor.CreateUser(newUser)
	if err != nil {
		log.Printf("failed create user, err: %s", err)
		helper.GenerateResponseAPI(http.StatusInternalServerError, "error", "Failed create user", c, false)
		return
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", user, c, false)
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.InputUserLogin

	if err := c.BindJSON(&input); err != nil {
		errsBinding := helper.ErrorBindingFormatter(err)
		log.Printf("failed to binding json input, err: %v", errsBinding)
		helper.GenerateResponseAPI(http.StatusBadRequest, "error binding", "Invalid input", c, false)
		return
	}

	u, err := h.userInteractor.GetUserByEmail(input.Email)
	if err != nil {
		log.Printf("failed to get user by email, email: %s, err: %s", input.Email, err)
		helper.GenerateResponseAPI(http.StatusNotFound, "error", "User not found", c, false)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		log.Printf("failed to get user by email, email: %s", input.Email)
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", "Email or password invalid", c, false)
		return
	}

	token, err := h.authService.GenerateToken(u.Id)
	if err != nil {
		log.Printf("failed to generate token, err: %s", err)
		helper.GenerateResponseAPI(http.StatusInternalServerError, "error generate token", "Failed to generate token", c, false)
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

	user, err := h.userInteractor.GetUserById(userId)
	if err != nil {
		helper.GenerateResponseAPI(0, "error", err.Error(), c, false)
		return
	}

	userWithPosts := response.UserWithPostsResponse{
		User: user,
		Post: posts,
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", userWithPosts, c, false)
}
