package mocks

import (
	"context"

	"github.com/arrase21/crm-users/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockUserRepo is a mock implementation of domain.UserRepo
type MockUserRepo struct {
	mock.Mock
}

// NewMockUserRepo creates a new instance of MockUserRepo
func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{}
}

// Create provides a mock function with given fields: ctx, usr
func (m *MockUserRepo) Create(ctx context.Context, usr *domain.User) error {
	args := m.Called(ctx, usr)
	return args.Error(0)
}

// GetByID provides a mock function with given fields: ctx, id
func (m *MockUserRepo) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// GetByDni provides a mock function with given fields: ctx, dni
func (m *MockUserRepo) GetByDni(ctx context.Context, dni string) (*domain.User, error) {
	args := m.Called(ctx, dni)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// List provides a mock function with given fields: ctx, page, limit
func (m *MockUserRepo) List(ctx context.Context, page, limit int) ([]domain.User, int64, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]domain.User), args.Get(1).(int64), args.Error(2)
}

// Update provides a mock function with given fields: ctx, usr
func (m *MockUserRepo) Update(ctx context.Context, usr *domain.User) error {
	args := m.Called(ctx, usr)
	return args.Error(0)
}

// Delete provides a mock function with given fields: ctx, id
func (m *MockUserRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
