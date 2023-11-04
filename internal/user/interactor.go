package user

import (
	"errors"
	"log"
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
	ErrorUserNotFound    = errors.New("user not found")
)

func NewInteractor(userRepo UserRepository) *interactor {
	return &interactor{
		userRepo: userRepo,
	}
}

func (i *interactor) CreateUser(user User) (User, error) {
	user, err := i.userRepo.Save(user)
	if err != nil {
		log.Printf("failed to create user, err: %s", err)
		return user, ErrorDatabaseFailure
	}

	return user, nil
}

func (i *interactor) ValidateUser(userId int) (bool, error) {
	_, err := i.userRepo.FindById(userId)
	if err != nil {
		if errors.Is(err, ErrorUserNotFound) {
			log.Printf("failed to find user with id: %d, err: %s", userId, err)
			return false, ErrorUserNotFound
		}

		log.Printf("failed to find user with id: %d, err: %s", userId, err)
		return false, ErrorDatabaseFailure
	}

	return true, nil
}

func (i *interactor) GetUserById(userId int) (User, error) {
	user, err := i.userRepo.FindById(userId)
	if err != nil {
		if errors.Is(err, ErrorUserNotFound) {
			log.Printf("failed to find user with id: %d, err: %s", userId, err)
			return user, ErrorUserNotFound
		}

		log.Printf("failed to find user with id: %d, err: %s", userId, err)
		return user, ErrorDatabaseFailure
	}
	return user, nil
}

func (i *interactor) GetUserByEmail(email string) (User, error) {
	user, err := i.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, ErrorUserNotFound) {
			log.Printf("failed to find user with email: %s, err: %s", email, err)
			return user, ErrorUserNotFound
		}

		log.Printf("failed to find user with email: %s, err: %s", email, err)
		return user, ErrorDatabaseFailure
	}
	return user, nil
}
