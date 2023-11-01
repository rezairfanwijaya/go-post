package user

import (
	"errors"
)

type Interactor interface {
	CreateUser(user User) (User, error)
	ValidateUser(userId int) (bool, error)
	GetUserById(userId int) (User, error)
	GetUserByEmail(email string) (User, error)
}

type interactor struct {
	userRepo UserRepository
}

var (
	ErrorDatabaseFailure = errors.New("unknown error occurred")
	ErrorAuth            = errors.New("error access denied")
)

func NewInteractor(userRepo UserRepository) *interactor {
	return &interactor{
		userRepo: userRepo,
	}
}

func (i *interactor) CreateUser(user User) (User, error) {
	user, err := i.userRepo.Save(user)
	if err != nil {
		return user, ErrorDatabaseFailure
	}

	return user, nil
}

func (i *interactor) ValidateUser(userId int) (bool, error) {
	user, err := i.userRepo.FindById(userId)
	if err != nil {
		return false, err
	}

	if user.Id == 0 {
		return false, ErrorAuth
	}

	return true, nil
}

func (i *interactor) GetUserById(userId int) (User, error) {
	user, err := i.userRepo.FindById(userId)
	if err != nil {
		return user, ErrorDatabaseFailure
	}

	return user, nil
}

func (i *interactor) GetUserByEmail(email string) (User, error) {
	user, err := i.userRepo.FindByEmail(email)
	if err != nil {
		return user, ErrorDatabaseFailure
	}

	return user, nil
}
