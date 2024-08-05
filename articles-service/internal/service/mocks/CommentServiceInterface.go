// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	models "articles-service/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// CommentServiceInterface is an autogenerated mock type for the CommentServiceInterface type
type CommentServiceInterface struct {
	mock.Mock
}

// Create provides a mock function with given fields: comment
func (_m *CommentServiceInterface) Create(comment *models.Comment) (int, error) {
	ret := _m.Called(comment)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Comment) (int, error)); ok {
		return rf(comment)
	}
	if rf, ok := ret.Get(0).(func(*models.Comment) int); ok {
		r0 = rf(comment)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(*models.Comment) error); ok {
		r1 = rf(comment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCommentsOnArticle provides a mock function with given fields: articleId
func (_m *CommentServiceInterface) GetCommentsOnArticle(articleId int) ([]*models.Comment, error) {
	ret := _m.Called(articleId)

	if len(ret) == 0 {
		panic("no return value specified for GetCommentsOnArticle")
	}

	var r0 []*models.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]*models.Comment, error)); ok {
		return rf(articleId)
	}
	if rf, ok := ret.Get(0).(func(int) []*models.Comment); ok {
		r0 = rf(articleId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Comment)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(articleId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCommentServiceInterface creates a new instance of CommentServiceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCommentServiceInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *CommentServiceInterface {
	mock := &CommentServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}