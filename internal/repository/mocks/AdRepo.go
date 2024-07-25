// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	models "ads-service/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// AdRepo is an autogenerated mock type for the AdRepo type
type AdRepo struct {
	mock.Mock
}

// Create provides a mock function with given fields: ad
func (_m *AdRepo) Create(ad *models.Ad) (int, error) {
	ret := _m.Called(ad)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Ad) (int, error)); ok {
		return rf(ad)
	}
	if rf, ok := ret.Get(0).(func(*models.Ad) int); ok {
		r0 = rf(ad)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(*models.Ad) error); ok {
		r1 = rf(ad)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: priceSort, dateSort, page
func (_m *AdRepo) GetAll(priceSort string, dateSort string, page int) ([]*models.Ad, error) {
	ret := _m.Called(priceSort, dateSort, page)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []*models.Ad
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, int) ([]*models.Ad, error)); ok {
		return rf(priceSort, dateSort, page)
	}
	if rf, ok := ret.Get(0).(func(string, string, int) []*models.Ad); ok {
		r0 = rf(priceSort, dateSort, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Ad)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, int) error); ok {
		r1 = rf(priceSort, dateSort, page)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOne provides a mock function with given fields: id
func (_m *AdRepo) GetOne(id int) (*models.Ad, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetOne")
	}

	var r0 *models.Ad
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*models.Ad, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *models.Ad); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Ad)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAdRepo creates a new instance of AdRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAdRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *AdRepo {
	mock := &AdRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
