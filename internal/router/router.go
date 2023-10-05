package router

import (
	"database/sql"
	"go-post/internal/article"
	"go-post/internal/post"
	"go-post/internal/user"

	"github.com/gin-gonic/gin"
)

func NewRouter(engine *gin.Engine, db *sql.DB) {
	userRepo := user.NewRepositoryUser(db)
	postRepo := post.NewPostRepository(db)
	userHandler := user.NewUserHandler(userRepo)
	postHandler := post.NewPostHandler(postRepo)
	articleHandler := article.NewArticleHandler(userRepo, postRepo)

	API := engine.Group("/api")
	API.POST("/user/signup", userHandler.SignUp)

	API.POST("/post", postHandler.CreatePost)
	API.DELETE("/posts/:id", postHandler.DeletePost)
	API.PUT("/posts/:id", postHandler.UpdatePost)

	API.GET("/articles/users/:id", articleHandler.GetArticlesByUser)
	API.GET("/articles/:id", articleHandler.GetArticlesDetail)
}
