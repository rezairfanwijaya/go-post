package repository

import (
	"database/sql"
	"go-post/internal/model"
)

type PostRepository interface {
	Save(post model.Post) error
	FindByPostId(postId int) (model.Post, error)
	FindByUserId(userId int) ([]model.Post, error)
	Delete(postId int) error
	Update(postId int, post model.Post) error
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

func (r *postRepository) FindByUserId(userId int) ([]model.Post, error) {
	posts := []model.Post{}

	query := `
		SELECT id, title, content FROM posts WHERE user_id = $1
	`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return posts, err
	}

	for rows.Next() {
		post := model.Post{}

		err := rows.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
		)

		if err != nil {
			return posts, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *postRepository) Delete(postId int) error {
	query := `
		DELETE FROM posts WHERE id = $1
	`

	_, err := r.db.Exec(query, postId)
	if err != nil {
		return err
	}

	return nil
}

func (r *postRepository) Update(postId int, post model.Post) error {
	query := `
		UPDATE posts SET title = $1, content = $2 WHERE id = $3
	`

	_, err := r.db.Exec(query, post.Title, post.Content, postId)

	if err != nil {
		return err
	}

	return nil
}
