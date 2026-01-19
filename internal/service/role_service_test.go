package service

import (
	"context"
	"errors"
	"testing"

	"github.com/arrase21/crm-users/internal/domain"
	"github.com/arrase21/crm-users/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoleService_Create(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockRoleRepo()
	service := NewRoleService(mockRepo)
	ctx := context.Background()

	// Test data
	validRole := &domain.Role{
		TenantID:    1,
		Name:        "admin",
		Description: "Administrator role",
		IsSystem:    false,
		IsActive:    true,
	}

	t.Run("✅ Success - Valid role creation", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByName", ctx, validRole.Name).Return(nil, domain.ErrRoleNotFound).Once()
		mockRepo.On("Create", ctx, validRole).Return(nil).Once()

		// Act
		err := service.Create(ctx, validRole)

		// Assert
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("❌ Error - Nil role", func(t *testing.T) {
		// Act
		err := service.Create(ctx, nil)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "role cannot be nil")
		mockRepo.AssertNotCalled(t, "GetByName")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("❌ Error - Empty role name", func(t *testing.T) {
		// Arrange
		invalidRole := *validRole
		invalidRole.Name = ""

		// Act
		err := service.Create(ctx, &invalidRole)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "role name is required")
		mockRepo.AssertNotCalled(t, "GetByName")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("❌ Error - Role name already exists", func(t *testing.T) {
		// Arrange
		existingRole := &domain.Role{ID: 1, Name: validRole.Name}
		mockRepo.On("GetByName", ctx, validRole.Name).Return(existingRole, nil).Once()

		// Act
		err := service.Create(ctx, validRole)

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrRoleExisting))
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("❌ Error - Repository error on name check", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByName", ctx, validRole.Name).Return(nil, errors.New("database error")).Once()

		// Act
		err := service.Create(ctx, validRole)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error checking existing role")
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("❌ Error - Repository error on create", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByName", ctx, validRole.Name).Return(nil, domain.ErrRoleNotFound).Once()
		mockRepo.On("Create", ctx, validRole).Return(errors.New("create error")).Once()

		// Act
		err := service.Create(ctx, validRole)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create error")
		mockRepo.AssertExpectations(t)
	})
}

func TestRoleService_GetByID(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockRoleRepo()
	service := NewRoleService(mockRepo)
	ctx := context.Background()

	t.Run("✅ Success - Valid role ID", func(t *testing.T) {
		// Arrange
		expectedRole := &domain.Role{
			ID:          1,
			Name:        "admin",
			Description: "Administrator role",
		}
		mockRepo.On("GetByID", ctx, uint(1)).Return(expectedRole, nil).Once()

		// Act
		role, err := service.GetByID(ctx, 1)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedRole, role)
		mockRepo.AssertExpectations(t)
	})

	t.Run("❌ Error - Invalid role ID (zero)", func(t *testing.T) {
		// Act
		role, err := service.GetByID(ctx, 0)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, role)
		assert.Contains(t, err.Error(), "invalid role id")
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("❌ Error - Role not found", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByID", ctx, uint(999)).Return(nil, domain.ErrRoleNotFound).Once()

		// Act
		role, err := service.GetByID(ctx, 999)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, role)
		assert.True(t, errors.Is(err, domain.ErrRoleNotFound))
		mockRepo.AssertExpectations(t)
	})
}

func TestRoleService_GetByName(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockRoleRepo()
	service := NewRoleService(mockRepo)
	ctx := context.Background()

	t.Run("✅ Success - Valid role name", func(t *testing.T) {
		// Arrange
		roleName := "admin"
		expectedRole := &domain.Role{
			ID:   1,
			Name: roleName,
		}
		mockRepo.On("GetByName", ctx, roleName).Return(expectedRole, nil).Once()

		// Act
		role, err := service.GetByName(ctx, roleName)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedRole, role)
		mockRepo.AssertExpectations(t)
	})

	t.Run("❌ Error - Empty role name", func(t *testing.T) {
		// Act
		role, err := service.GetByName(ctx, "")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, role)
		assert.Contains(t, err.Error(), "invalid name")
		mockRepo.AssertNotCalled(t, "GetByName")
	})

	t.Run("❌ Error - Role not found", func(t *testing.T) {
		// Arrange
		roleName := "nonexistent"
		mockRepo.On("GetByName", ctx, roleName).Return(nil, domain.ErrRoleNotFound).Once()

		// Act
		role, err := service.GetByName(ctx, roleName)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, role)
		assert.True(t, errors.Is(err, domain.ErrRoleNotFound))
		mockRepo.AssertExpectations(t)
	})
}

func TestRoleService_List(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockRoleRepo()
	service := NewRoleService(mockRepo)
	ctx := context.Background()

	t.Run("✅ Success - List roles", func(t *testing.T) {
		// Arrange
		expectedRoles := []domain.Role{
			{ID: 1, Name: "admin", Description: "Administrator role"},
			{ID: 2, Name: "user", Description: "Regular user role"},
		}
		mockRepo.On("List", ctx).Return(expectedRoles, nil).Once()

		// Act
		roles, err := service.List(ctx)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedRoles, roles)
		mockRepo.AssertExpectations(t)
	})

	t.Run("❌ Error - Repository error", func(t *testing.T) {
		// Arrange
		mockRepo.On("List", ctx).Return(nil, errors.New("database error")).Once()

		// Act
		roles, err := service.List(ctx)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, roles)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestRoleService_Update(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockRoleRepo()
	service := NewRoleService(mockRepo)
	ctx := context.Background()

	// Test data
	validRole := &domain.Role{
		ID:          1,
		TenantID:    1,
		Name:        "admin",
		Description: "Updated administrator role",
		IsSystem:    false,
		IsActive:    true,
	}

	t.Run("✅ Success - Valid role update", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByID", ctx, validRole.ID).Return(validRole, nil).Once()
		mockRepo.On("GetByName", ctx, validRole.Name).Return(nil, domain.ErrRoleNotFound).Once()
		mockRepo.On("Update", ctx, validRole).Return(nil).Once()

		// Act
		err := service.Update(ctx, validRole)

		// Assert
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("✅ Success - Update with same name", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByID", ctx, validRole.ID).Return(validRole, nil).Once()
		mockRepo.On("GetByName", ctx, validRole.Name).Return(validRole, nil).Once()
		mockRepo.On("Update", ctx, validRole).Return(nil).Once()

		// Act
		err := service.Update(ctx, validRole)

		// Assert
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("❌ Error - Nil role", func(t *testing.T) {
		// Act
		err := service.Update(ctx, nil)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid role")
		mockRepo.AssertNotCalled(t, "GetByID")
		mockRepo.AssertNotCalled(t, "GetByName")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("❌ Error - Missing role ID", func(t *testing.T) {
		// Arrange
		roleWithoutID := *validRole
		roleWithoutID.ID = 0

		// Act
		err := service.Update(ctx, &roleWithoutID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid role")
		mockRepo.AssertNotCalled(t, "GetByID")
		mockRepo.AssertNotCalled(t, "GetByName")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("❌ Error - Role not found by ID", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByID", ctx, validRole.ID).Return(nil, domain.ErrRoleNotFound).Once()

		// Act
		err := service.Update(ctx, validRole)

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrRoleNotFound))
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "GetByName")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("❌ Error - Role name already exists for another role", func(t *testing.T) {
		// Arrange
		anotherRole := &domain.Role{ID: 2, Name: validRole.Name}
		mockRepo.On("GetByID", ctx, validRole.ID).Return(validRole, nil).Once()
		mockRepo.On("GetByName", ctx, validRole.Name).Return(anotherRole, nil).Once()

		// Act
		err := service.Update(ctx, validRole)

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrRoleExisting))
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Update")
	})
}

func TestRoleService_Delete(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockRoleRepo()
	service := NewRoleService(mockRepo)
	ctx := context.Background()

	t.Run("✅ Success - Valid role deletion", func(t *testing.T) {
		// Arrange
		roleID := uint(1)
		mockRepo.On("Delete", ctx, roleID).Return(nil).Once()

		// Act
		err := service.Delete(ctx, roleID)

		// Assert
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("❌ Error - Invalid role ID (zero)", func(t *testing.T) {
		// Act
		err := service.Delete(ctx, 0)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid role id")
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("❌ Error - Repository error", func(t *testing.T) {
		// Arrange
		roleID := uint(1)
		mockRepo.On("Delete", ctx, roleID).Return(errors.New("database error")).Once()

		// Act
		err := service.Delete(ctx, roleID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}
