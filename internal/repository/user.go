package repository

import (
	"go-post/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	Save(user model.User) (model.User, error)
	FindByEmail(email string) (model.User, error)
	FindById(userId string) (model.User, error)
}

type repository struct {
	db *gorm.DB
}

func (r *repository) Save(user model.User) (model.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByEmail(email string) (model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindById(userId string) (model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", userId).Error; err != nil {
		return user, err
	}

	return user, nil
}
