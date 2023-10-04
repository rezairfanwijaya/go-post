package repository

import (
	"database/sql"
	"go-post/internal/model"
)

type PostRepository interface {
	Save(post model.Post) error
	FindByPostId(postId int) (model.Post, error)
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

func (r *postRepository) FindByPostId(postId int) (model.Post, error) {
	post := model.Post{}

	query := `
		SELECT * FROM posts WHERE id = $1
	`

	row := r.db.QueryRow(query, postId)
	err := row.Scan(
		&post.Id,
		&post.UserId,
		&post.Title,
		&post.Content,
	)

	if err != nil {
		return post, err
	}

	return post, nil
}
