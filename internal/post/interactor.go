package post

import (
	"fmt"
	"net/http"
)

type Interactor interface {
	CreatePost(post Post) (int, error)
	GetPost(userId, postId int) (Post, int, error)
	GetPostByUserId(userId int) ([]Post, int, error)
	UpdatePost(postId int, input InputUpdatePost) (Post, int, error)
	DeletePost(postId, userId int) (int, error)
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

	if isValid := validateUser(userId, post); !isValid {
		return post, http.StatusUnauthorized, fmt.Errorf("unauthorized")
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

func (i *interactor) UpdatePost(postId int, input InputUpdatePost) (Post, int, error) {
	post, err := i.postRepository.FindByPostId(postId)
	if err != nil {
		return post, http.StatusInternalServerError, err
	}

	if isValid := validateUser(input.UserId, post); !isValid {
		return post, http.StatusUnauthorized, fmt.Errorf("unauthorized")
	}

	post.Content = input.Content
	post.Title = input.Title

	err = i.postRepository.Update(postId, post)
	if err != nil {
		return post, http.StatusInternalServerError, err
	}

	return post, http.StatusOK, nil
}

func (i *interactor) DeletePost(postId, userId int) (int, error) {
	post, err := i.postRepository.FindByPostId(postId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if isValid := validateUser(userId, post); !isValid {
		return http.StatusUnauthorized, fmt.Errorf("unauthorized")
	}

	err = i.postRepository.Delete(post.Id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func validateUser(userId int, post Post) bool {
	return userId == post.UserId
}
