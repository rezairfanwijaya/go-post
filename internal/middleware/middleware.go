package middleware

import (
	"go-post/internal/auth"
	"go-post/internal/helper"
	"go-post/internal/user"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

func AuthFunc(authService auth.AuthService, userInteractor user.Interactor) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			helper.GenerateResponseAPI(http.StatusBadRequest, "error", "invalid token format", c, true)
			return
		}

		tokenSplit := strings.Split(authHeader, " ")

		token := ""
		if len(tokenSplit) == 2 {
			token = tokenSplit[1]
		}

		tokenVerified, err := authService.VerifyToken(token)
		if err != nil {
			helper.GenerateResponseAPI(http.StatusUnauthorized, "unauthorized", "unauthorized", c, true)
			return
		}

		claim, ok := tokenVerified.Claims.(jwt.MapClaims)
		if !ok {
			helper.GenerateResponseAPI(http.StatusUnauthorized, "unauthorized", "unauthorized", c, true)
			return
		}

		userId := int(claim["user_id"].(float64))

		isValid, err := userInteractor.ValidateUser(userId)
		if !isValid {
			helper.GenerateResponseAPI(http.StatusUnauthorized, "unauthorized", err.Error(), c, true)
			return
		}

		c.Set("user", userId)
	}
}
