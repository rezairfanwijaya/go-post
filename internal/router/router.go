package router

import (
	"database/sql"
	"go-post/internal/auth"
	"go-post/internal/handler"
	"go-post/internal/middleware"
	"go-post/internal/post"
	"go-post/internal/user"

	"github.com/gin-gonic/gin"
)

func NewRouter(engine *gin.Engine, db *sql.DB) {
	authService := auth.NewAuthService()

	userRepo := user.NewRepositoryUser(db)
	postRepo := post.NewPostRepository(db)

	userInteractor := user.NewInteractor(userRepo)
	postInteractor := post.NewInteractor(postRepo)

	userHandler := handler.NewUserHandler(postInteractor, userInteractor, authService)
	postHandler := handler.NewPostHandler(userInteractor, postInteractor)

	API := engine.Group("/api")

	API.POST("/users/signup", userHandler.SignUp)
	API.POST("/users/login", userHandler.Login)
	API.GET("/users/posts", middleware.AuthFunc(authService, userInteractor), userHandler.GetUserWithPosts)

	API.POST("/posts", middleware.AuthFunc(authService, userInteractor), postHandler.CreatePost)
	API.DELETE("/posts/:id", middleware.AuthFunc(authService, userInteractor), postHandler.DeletePost)
	API.PUT("/posts/:id", middleware.AuthFunc(authService, userInteractor), postHandler.UpdatePost)
	API.GET("/posts/:id", middleware.AuthFunc(authService, userInteractor), postHandler.GetPost)
}
