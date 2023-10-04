package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewConnection() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"127.0.0.1",
		"5432",
		"postgres",
		"admin",
		"article",
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return db, err
	}

	migrationSchemaUser := `
		CREATE TABLE if not exists users (
			id SERIAL primary key,
			email varchar(200),
			password varchar(200),
			CONSTRAINT email_unique UNIQUE (email)
		)
	`

	migrationSchemaPost := `
		CREATE TABLE if not exists posts (
			id SERIAL primary key,
			user_id int,
			title varchar(200),
			content varchar(500),
			CONSTRAINT title_unique UNIQUE (title)
		)
	`

	_, err = db.Exec(migrationSchemaUser)
	if err != nil {
		return db, err
	}

	_, err = db.Exec(migrationSchemaPost)
	if err != nil {
		return db, err
	}

	return db, nil
}
