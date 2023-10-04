package middleware

import (
	"go-post/internal/auth"
	"go-post/internal/helper"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(auth auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			helper.GenerateResponseAPI(http.StatusUnauthorized, "Unauthorized", nil, c, true)
			return
		}

		var tokenHeader string
		headerSplit := strings.Split(authHeader, " ")
		if len(headerSplit) == 2 {
			tokenHeader = headerSplit[1]
		}

		token, err := auth.ValidateToken(tokenHeader)
		if err != nil {
			helper.GenerateResponseAPI(http.StatusUnauthorized, "Unauthorized", nil, c, true)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || token.Valid {
			helper.GenerateResponseAPI(http.StatusUnauthorized, "Unauthorized", nil, c, true)
			return
		}

		userID := claims["user_id"].(int)

		c.Set("currentUser", userID)
	}
}
