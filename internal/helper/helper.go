package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type responseAPI struct {
	Meta meta `json:"meta"`
	Data any  `json:"data"`
}

type meta struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func GenerateResponseAPI(code int, status string, data any, c *gin.Context, isMiddleware bool) {
	response := responseAPI{
		Meta: meta{
			Code:   code,
			Status: status,
		},
		Data: data,
	}

	if isMiddleware {
		c.AbortWithStatusJSON(code, response)
	} else {
		c.JSON(code, response)
	}
}

func ErrorBindingFormatter(err error) []string {
	var errBindings []string
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, e := range err.(validator.ValidationErrors) {
			errBindings = append(errBindings, e.Error())
		}
	}

	return errBindings
}
