package main

import (
	"go-post/internal/auth"
	"go-post/internal/database"
	"go-post/internal/handler"
	"go-post/internal/repository"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal(err)
	}

	authService := auth.NewAuthService()

	userRepo := repository.NewRepositoryUser(db)
	userHandler := handler.NewHandlerUser(userRepo, authService)

	r := gin.Default()

	API := r.Group("/api")
	API.POST("/user/signup", userHandler.SignUp)
	API.POST("/user/login", userHandler.Login)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
