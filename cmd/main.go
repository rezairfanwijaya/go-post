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
	postRepo := repository.NewPostRepository(db)
	userHandler := handler.NewUserHandler(userRepo, postRepo)
	postHandler := handler.NewPostHandler(postRepo, userRepo)

	r := gin.Default()

	API := r.Group("/api")
	API.POST("/user/signup", userHandler.SignUp)
	API.POST("/user/Create", userHandler.Login)
	API.GET("/users/:id/posts", userHandler.GetUserWitPosts)

	API.POST("/post", postHandler.CreatePost)
	API.GET("/posts/:id/user", postHandler.GetPostDetail)
	API.DELETE("/posts/:id", postHandler.DeletePost)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
