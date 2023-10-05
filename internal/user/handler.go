package user

import (
	"go-post/internal/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type userHandler struct {
	userRepo UserRepository
}

func NewUserHandler(userRepo UserRepository) *userHandler {
	return &userHandler{
		userRepo: userRepo,
	}
}

func (h *userHandler) SignUp(c *gin.Context) {
	var input InputUserSignUp

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

	var newUser User
	newUser.Email = input.Email
	newUser.Password = string(passwordHashed)

	err = h.userRepo.Save(newUser)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", err, c, false)
		return
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", "success", c, false)
}
