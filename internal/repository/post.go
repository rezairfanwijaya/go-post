package repository

import (
	"go-post/internal/model"

	"gorm.io/gorm"
)

type PostRepository interface {
	Save(post model.Post) (model.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) Save(post model.Post) (model.Post, error) {
	if err := r.db.Create(&post).Error; err != nil {
		return post, err
	}

	return post, nil
}
