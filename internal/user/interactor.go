package user

import "net/http"

type Interactor interface {
	CreateUser(user User) (int, error)
	ValidateUser(userId int) (bool, error)
	GetUserById(userId int) (User, int, error)
	GetUserByEmail(email string) (User, int, error)
}

type interactor struct {
	userRepo UserRepository
}

func NewInteractor(userRepo UserRepository) *interactor {
	return &interactor{
		userRepo: userRepo,
	}
}

func (i *interactor) CreateUser(user User) (int, error) {
	if err := i.userRepo.Save(user); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (i *interactor) ValidateUser(userId int) (bool, error) {
	user, err := i.userRepo.FindById(userId)
	if err != nil {
		return false, err
	}

	return user.Id != 0, nil
}

func (i *interactor) GetUserById(userId int) (User, int, error) {
	user, err := i.userRepo.FindById(userId)
	if err != nil {
		return user, http.StatusInternalServerError, err
	}

	return user, http.StatusOK, nil
}

func (i *interactor) GetUserByEmail(email string) (User, int, error) {
	user, err := i.userRepo.FindByEmail(email)
	if err != nil {
		return user, http.StatusInternalServerError, err
	}

	return user, http.StatusOK, nil
}
