package database

import (
	"database/sql"
	"fmt"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func NewConnection(path string) (*sql.DB, error) {
	env, err := godotenv.Read(path)
	if err != nil {
		return &sql.DB{}, err
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s  dbname=%s sslmode=disable",
		env["HOST"],
		env["PORT"],
		env["USER"],
		env["DBNAME"],
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
