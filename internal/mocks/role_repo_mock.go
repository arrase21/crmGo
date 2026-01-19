package mocks

import (
	"context"

	"github.com/arrase21/crm-users/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockRoleRepo is a mock implementation of domain.RoleRepo
type MockRoleRepo struct {
	mock.Mock
}

// NewMockRoleRepo creates a new instance of MockRoleRepo
func NewMockRoleRepo() *MockRoleRepo {
	return &MockRoleRepo{}
}

// Create provides a mock function with given fields: ctx, role
func (m *MockRoleRepo) Create(ctx context.Context, role *domain.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

// GetByID provides a mock function with given fields: ctx, roleID
func (m *MockRoleRepo) GetByID(ctx context.Context, roleID uint) (*domain.Role, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

// GetByName provides a mock function with given fields: ctx, name
func (m *MockRoleRepo) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

// List provides a mock function with given fields: ctx
func (m *MockRoleRepo) List(ctx context.Context) ([]domain.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Role), args.Error(1)
}

// Update provides a mock function with given fields: ctx, role
func (m *MockRoleRepo) Update(ctx context.Context, role *domain.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

// Delete provides a mock function with given fields: ctx, roleID
func (m *MockRoleRepo) Delete(ctx context.Context, roleID uint) error {
	args := m.Called(ctx, roleID)
	return args.Error(0)
}

// AssignPermission provides a mock function with given fields: ctx, roleID, actionID
func (m *MockRoleRepo) AssignPermission(ctx context.Context, roleID, actionID uint) error {
	args := m.Called(ctx, roleID, actionID)
	return args.Error(0)
}

// RevokePermission provides a mock function with given fields: ctx, roleID, actionID
func (m *MockRoleRepo) RevokePermission(ctx context.Context, roleID, actionID uint) error {
	args := m.Called(ctx, roleID, actionID)
	return args.Error(0)
}

// GetPermissions provides a mock function with given fields: ctx, roleID
func (m *MockRoleRepo) GetPermissions(ctx context.Context, roleID uint) ([]domain.PermissionAction, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PermissionAction), args.Error(1)
}
