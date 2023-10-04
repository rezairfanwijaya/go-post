package main

import (
	"go-post/internal/database"
	"log"
)

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(db)
}
