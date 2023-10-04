package repository

import (
	"go-post/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Save(user model.User) (model.User, error)
	FindByEmail(email string) (model.User, error)
	FindById(userId string) (model.User, error)
}

type serRepository struct {
	db *gorm.DB
}

func NewRepositoryUser(db *gorm.DB) UserRepository {
	return &serRepository{
		db: db,
	}
}

func (r *serRepository) Save(user model.User) (model.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *serRepository) FindByEmail(email string) (model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).Find(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *serRepository) FindById(userId string) (model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", userId).Error; err != nil {
		return user, err
	}

	return user, nil
}
