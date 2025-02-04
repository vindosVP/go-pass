// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	models "github.com/vindosVP/go-pass/internal/models"
)

// PasswordStorage is an autogenerated mock type for the PasswordStorage type
type PasswordStorage struct {
	mock.Mock
}

// AddPassword provides a mock function with given fields: ctx, pwd
func (_m *PasswordStorage) AddPassword(ctx context.Context, pwd *models.Password) (int, error) {
	ret := _m.Called(ctx, pwd)

	if len(ret) == 0 {
		panic("no return value specified for AddPassword")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Password) (int, error)); ok {
		return rf(ctx, pwd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Password) int); ok {
		r0 = rf(ctx, pwd)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Password) error); ok {
		r1 = rf(ctx, pwd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeletePassword provides a mock function with given fields: ctx, id, ownerID
func (_m *PasswordStorage) DeletePassword(ctx context.Context, id int, ownerID int) error {
	ret := _m.Called(ctx, id, ownerID)

	if len(ret) == 0 {
		panic("no return value specified for DeletePassword")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) error); ok {
		r0 = rf(ctx, id, ownerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetPasswords provides a mock function with given fields: ctx, ownerID
func (_m *PasswordStorage) GetPasswords(ctx context.Context, ownerID int) ([]*models.Password, error) {
	ret := _m.Called(ctx, ownerID)

	if len(ret) == 0 {
		panic("no return value specified for GetPasswords")
	}

	var r0 []*models.Password
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) ([]*models.Password, error)); ok {
		return rf(ctx, ownerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) []*models.Password); ok {
		r0 = rf(ctx, ownerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Password)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, ownerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePassword provides a mock function with given fields: ctx, pwd
func (_m *PasswordStorage) UpdatePassword(ctx context.Context, pwd *models.Password) error {
	ret := _m.Called(ctx, pwd)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePassword")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Password) error); ok {
		r0 = rf(ctx, pwd)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPasswordStorage creates a new instance of PasswordStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPasswordStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *PasswordStorage {
	mock := &PasswordStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
