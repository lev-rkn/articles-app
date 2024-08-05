// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	models "articles-service/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// ArticleServiceInterface is an autogenerated mock type for the ArticleServiceInterface type
type ArticleServiceInterface struct {
	mock.Mock
}

// Create provides a mock function with given fields: ad
func (_m *ArticleServiceInterface) Create(ad *models.Article) (int, error) {
	ret := _m.Called(ad)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Article) (int, error)); ok {
		return rf(ad)
	}
	if rf, ok := ret.Get(0).(func(*models.Article) int); ok {
		r0 = rf(ad)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(*models.Article) error); ok {
		r1 = rf(ad)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: priceSort, dateSort, page, userId
func (_m *ArticleServiceInterface) GetAll(priceSort string, dateSort string, page int, userId int) ([]*models.Article, error) {
	ret := _m.Called(priceSort, dateSort, page, userId)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []*models.Article
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, int, int) ([]*models.Article, error)); ok {
		return rf(priceSort, dateSort, page, userId)
	}
	if rf, ok := ret.Get(0).(func(string, string, int, int) []*models.Article); ok {
		r0 = rf(priceSort, dateSort, page, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Article)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, int, int) error); ok {
		r1 = rf(priceSort, dateSort, page, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOne provides a mock function with given fields: id
func (_m *ArticleServiceInterface) GetOne(id int) (*models.Article, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetOne")
	}

	var r0 *models.Article
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*models.Article, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *models.Article); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Article)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewArticleServiceInterface creates a new instance of ArticleServiceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewArticleServiceInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *ArticleServiceInterface {
	mock := &ArticleServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}