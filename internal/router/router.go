package router

import (
	"database/sql"
	"go-post/internal/handler"
	"go-post/internal/post"
	"go-post/internal/user"

	"github.com/gin-gonic/gin"
)

func NewRouter(engine *gin.Engine, db *sql.DB) {
	userRepo := user.NewRepositoryUser(db)
	postRepo := post.NewPostRepository(db)
	userHandler := handler.NewUserHandler(userRepo, postRepo)
	postHandler := handler.NewPostHandler(postRepo, userRepo)

	API := engine.Group("/api")
	API.POST("/users/signup", userHandler.SignUp)
	API.GET("/users/:id/posts", userHandler.GetUserWithPosts)

	API.POST("/posts", postHandler.CreatePost)
	API.DELETE("/posts/:id", postHandler.DeletePost)
	API.PUT("/posts/:id", postHandler.UpdatePost)
	API.GET("/posts/:id", postHandler.GetPost)
}
