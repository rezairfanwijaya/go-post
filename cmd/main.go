package main

import (
	"go-post/internal/database"
	"go-post/internal/router"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	router.NewRouter(r, db)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
