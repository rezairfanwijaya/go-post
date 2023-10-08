package post

import (
	"fmt"
	"net/http"
)

type Interactor interface {
	CreatePost(post Post) (int, error)
	GetPost(userId, postId int) (Post, int, error)
	GetPostByUserId(userId int) ([]Post, int, error)
	UpdatePost(postId, userId int, post Post) (Post, int, error)
	DeletePost(postId, userId int) (int, error)
	ValidateUser(userId int, post Post) bool
}

type interactor struct {
	postRepository PostRepository
}

func NewInteractor(postRepository PostRepository) Interactor {
	return &interactor{
		postRepository: postRepository,
	}
}

func (i *interactor) CreatePost(post Post) (int, error) {
	if err := i.postRepository.Save(post); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (i *interactor) GetPost(userId, postId int) (Post, int, error) {
	post, err := i.postRepository.FindByPostId(postId)
	if err != nil {
		return post, http.StatusInternalServerError, err
	}

	if isValid := i.ValidateUser(userId, post); !isValid {
		return Post{}, http.StatusUnauthorized, fmt.Errorf("unauthorized")
	}

	return post, http.StatusOK, nil
}

func (i *interactor) GetPostByUserId(userId int) ([]Post, int, error) {
	posts, err := i.postRepository.FindByUserId(userId)
	if err != nil {
		return posts, http.StatusInternalServerError, err
	}

	return posts, http.StatusOK, nil
}

func (i *interactor) UpdatePost(postId, userId int, post Post) (Post, int, error) {
	if isValid := i.ValidateUser(userId, post); !isValid {
		return Post{}, http.StatusUnauthorized, fmt.Errorf("unauthorized")
	}

	err := i.postRepository.Update(postId, post)
	if err != nil {
		return Post{}, http.StatusInternalServerError, err
	}

	return post, http.StatusOK, nil
}

func (i *interactor) DeletePost(postId, userId int) (int, error) {
	post, err := i.postRepository.FindByPostId(postId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if isValid := i.ValidateUser(userId, post); !isValid {
		return http.StatusUnauthorized, fmt.Errorf("unauthorized")
	}

	err = i.postRepository.Delete(post.Id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (i *interactor) ValidateUser(userId int, post Post) bool {
	return userId == post.UserId
}
