package main

import (
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

	userRepo := repository.NewRepositoryUser(db)
	userHandler := handler.NewHandlerUser(userRepo)

	r := gin.Default()

	API := r.Group("/api")
	API.POST("/user/signup", userHandler.SignUp)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
