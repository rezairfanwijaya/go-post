package post

import (
	"database/sql"
)

type PostRepository interface {
	Save(post Post) (Post, error)
	FindByPostId(postId int) (Post, error)
	FindByUserId(userId int) ([]Post, error)
	Delete(postId int) error
	Update(postId int, post Post) error
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) Save(post Post) (Post, error) {
	res := Post{}
	query := `
		INSERT into posts (user_id, title, content) VALUES ($1, $2, $3) RETURNING *
	`

	row := r.db.QueryRow(query, post.UserId, post.Title, post.Content)

	if err := row.Scan(&res.Id, &res.UserId, &res.Title, &res.Content); err != nil {
		return res, err
	}

	return res, nil
}

func (r *postRepository) FindByPostId(postId int) (Post, error) {
	post := Post{}

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

func (r *postRepository) FindByUserId(userId int) ([]Post, error) {
	posts := []Post{}

	query := `
		SELECT id, title, content FROM posts WHERE user_id = $1
	`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return posts, err
	}

	for rows.Next() {
		post := Post{}

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

func (r *postRepository) Update(postId int, post Post) error {
	query := `
		UPDATE posts SET title = $1, content = $2 WHERE id = $3
	`

	_, err := r.db.Exec(query, post.Title, post.Content, postId)

	if err != nil {
		return err
	}

	return nil
}
