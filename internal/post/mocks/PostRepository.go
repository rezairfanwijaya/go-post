// Code generated by mockery v2.35.2. DO NOT EDIT.

package mocks

import (
	post "go-post/internal/post"

	mock "github.com/stretchr/testify/mock"
)

// PostRepository is an autogenerated mock type for the PostRepository type
type PostRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: postId
func (_m *PostRepository) Delete(postId int) error {
	ret := _m.Called(postId)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(postId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByPostId provides a mock function with given fields: postId
func (_m *PostRepository) FindByPostId(postId int) (post.Post, error) {
	ret := _m.Called(postId)

	var r0 post.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (post.Post, error)); ok {
		return rf(postId)
	}
	if rf, ok := ret.Get(0).(func(int) post.Post); ok {
		r0 = rf(postId)
	} else {
		r0 = ret.Get(0).(post.Post)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(postId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByUserId provides a mock function with given fields: userId
func (_m *PostRepository) FindByUserId(userId int) ([]post.Post, error) {
	ret := _m.Called(userId)

	var r0 []post.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]post.Post, error)); ok {
		return rf(userId)
	}
	if rf, ok := ret.Get(0).(func(int) []post.Post); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]post.Post)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: _a0
func (_m *PostRepository) Save(_a0 post.Post) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(post.Post) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: postId, _a1
func (_m *PostRepository) Update(postId int, _a1 post.Post) error {
	ret := _m.Called(postId, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, post.Post) error); ok {
		r0 = rf(postId, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPostRepository creates a new instance of PostRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPostRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *PostRepository {
	mock := &PostRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
