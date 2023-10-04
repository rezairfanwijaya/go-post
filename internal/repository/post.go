package repository

import (
	"database/sql"
	"go-post/internal/model"
)

type PostRepository interface {
	Save(post model.Post) error
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) Save(post model.Post) error {
	query := `
		INSERT into posts (user_id, title, content) VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(query, post.UserId, post.Title, post.Content)
	if err != nil {
		return err
	}

	return nil

}
