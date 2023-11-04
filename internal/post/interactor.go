package post

import (
	"errors"
	"fmt"
	"log"
)

type Interactor interface {
	CreatePost(post Post) (Post, error)
	GetPost(userId, postId int) (Post, error)
	GetPostByUserId(userId int) ([]Post, error)
	UpdatePost(postId, userId int, post Post) (Post, error)
	DeletePost(postId, userId int) error
	ValidateUser(userId int, post Post) bool
}

type interactor struct {
	postRepository PostRepository
}

var (
	ErrDatabaseFailure = errors.New("unknown error occurred")
	ErrorPostNotFound  = errors.New("post not found")
	ErrorUnauthorized  = errors.New("unauthorized")
)

func NewInteractor(postRepository PostRepository) Interactor {
	return &interactor{
		postRepository: postRepository,
	}
}

func (i *interactor) CreatePost(post Post) (Post, error) {
	post, err := i.postRepository.Save(post)
	if err != nil {
		log.Printf("failed to save post, title: %s, userId: %d, err: %s", post.Title, post.UserId, err)
		return post, ErrDatabaseFailure
	}

	return post, nil
}

func (i *interactor) GetPost(userId, postId int) (Post, error) {
	post, err := i.postRepository.FindByPostId(postId)
	if err != nil {
		if errors.Is(err, ErrorPostNotFound) {
			log.Printf("failed to find post with title: %s, userId: %d, err: %s", post.Title, post.UserId, err)
			return post, ErrorPostNotFound
		}

		log.Printf("failed to find post with title: %s, userId: %d, err: %s", post.Title, post.UserId, err)
		return post, err
	}

	if isValid := i.ValidateUser(userId, post); !isValid {
		log.Printf("failed validate user, userId: %d, valid: %v", userId, isValid)
		return Post{}, ErrorUnauthorized
	}

	return post, nil
}

func (i *interactor) GetPostByUserId(userId int) ([]Post, error) {
	posts, err := i.postRepository.FindByUserId(userId)
	if err != nil {
		if errors.Is(err, ErrorPostNotFound) {
			log.Printf("failed to find post with userId: %d, err: %s", userId, err)
			return posts, ErrorPostNotFound
		}

		log.Printf("failed to find post with userId: %d, err: %s", userId, err)
		return posts, err
	}

	return posts, nil
}

func (i *interactor) UpdatePost(postId, userId int, post Post) (Post, error) {
	if isValid := i.ValidateUser(userId, post); !isValid {
		log.Printf("failed validate user, userId: %d, valid: %v", userId, isValid)
		return Post{}, fmt.Errorf("unauthorized")
	}

	err := i.postRepository.Update(postId, post)
	if err != nil {
		log.Printf("failed update post, userId: %d, postId: %d, err: %s", userId, postId, err)
		return Post{}, err
	}

	return post, nil
}

func (i *interactor) DeletePost(postId, userId int) error {
	post, err := i.postRepository.FindByPostId(postId)
	if err != nil {
		if errors.Is(err, ErrorPostNotFound) {
			log.Printf("failed to find post with userId: %d, err: %s", userId, err)
			return ErrorPostNotFound
		}

		log.Printf("failed to find post with userId: %d, err: %s", userId, err)
		return err
	}

	if isValid := i.ValidateUser(userId, post); !isValid {
		log.Printf("failed validate user, userId: %d, valid: %v", userId, isValid)
		return ErrorUnauthorized
	}

	err = i.postRepository.Delete(post.Id)
	if err != nil {
		log.Printf("failed delete post, err: %s", err)
		return err
	}

	return nil
}

func (i *interactor) ValidateUser(userId int, post Post) bool {
	return userId == post.UserId
}
