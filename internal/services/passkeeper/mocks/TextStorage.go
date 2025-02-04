// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	models "github.com/vindosVP/go-pass/internal/models"
)

// TextStorage is an autogenerated mock type for the TextStorage type
type TextStorage struct {
	mock.Mock
}

// AddText provides a mock function with given fields: ctx, t
func (_m *TextStorage) AddText(ctx context.Context, t *models.Text) (int, error) {
	ret := _m.Called(ctx, t)

	if len(ret) == 0 {
		panic("no return value specified for AddText")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Text) (int, error)); ok {
		return rf(ctx, t)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Text) int); ok {
		r0 = rf(ctx, t)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Text) error); ok {
		r1 = rf(ctx, t)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteText provides a mock function with given fields: ctx, id, ownerID
func (_m *TextStorage) DeleteText(ctx context.Context, id int, ownerID int) error {
	ret := _m.Called(ctx, id, ownerID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteText")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) error); ok {
		r0 = rf(ctx, id, ownerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetTexts provides a mock function with given fields: ctx, ownerID
func (_m *TextStorage) GetTexts(ctx context.Context, ownerID int) ([]*models.Text, error) {
	ret := _m.Called(ctx, ownerID)

	if len(ret) == 0 {
		panic("no return value specified for GetTexts")
	}

	var r0 []*models.Text
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) ([]*models.Text, error)); ok {
		return rf(ctx, ownerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) []*models.Text); ok {
		r0 = rf(ctx, ownerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Text)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, ownerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateText provides a mock function with given fields: ctx, t
func (_m *TextStorage) UpdateText(ctx context.Context, t *models.Text) error {
	ret := _m.Called(ctx, t)

	if len(ret) == 0 {
		panic("no return value specified for UpdateText")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Text) error); ok {
		r0 = rf(ctx, t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTextStorage creates a new instance of TextStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTextStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *TextStorage {
	mock := &TextStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
