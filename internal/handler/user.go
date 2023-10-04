package handler

import (
	"go-post/internal/helper"
	"go-post/internal/model"
	"go-post/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type handlerUser struct {
	repoUser repository.UserRepository
}

func NewHandlerUser(repoUser repository.UserRepository) *handlerUser {
	return &handlerUser{
		repoUser: repoUser,
	}
}

func (h *handlerUser) SignUp(c *gin.Context) {
	var input model.InputUserSignUp

	if err := c.BindJSON(&input); err != nil {
		errsBinding := helper.ErrorBindingFormatter(err)
		helper.GenerateResponseAPI(http.StatusBadRequest, "error binding", errsBinding, c)
		return
	}

	user, err := h.repoUser.FindByEmail(input.Email)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", err, c)
		return
	}

	if user.Id != 0 {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", "email already taken", c)
		return
	}

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error", err, c)
		return
	}

	var newUser model.User
	newUser.Email = input.Email
	newUser.Password = string(passwordHashed)

	userSaved, err := h.repoUser.Save(newUser)
	if err != nil {
		helper.GenerateResponseAPI(http.StatusBadRequest, "error binding", err, c)
		return
	}

	helper.GenerateResponseAPI(http.StatusOK, "success", userSaved, c)
}
