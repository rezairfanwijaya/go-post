package main

import (
	"go-post/internal/database"
	"go-post/internal/handler"
	"go-post/internal/repository"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewRepositoryUser(db)
	userHandler := handler.NewUserHandler(userRepo)
	postRepo := repository.NewPostRepository(db)
	postHandler := handler.NewPostHandler(postRepo)

	r := gin.Default()

	API := r.Group("/api")
	API.POST("/user/signup", userHandler.SignUp)
	API.POST("/user/Create", userHandler.Login)
	API.POST("/post", postHandler.NewPost)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
