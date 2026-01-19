package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
	"github.com/arrase21/crm-users/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_Update(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepo()
	service := NewUserService(mockRepo)
	ctx := context.Background()

	// Test data
	validUser := &domain.User{
		ID:        1,
		TenantID:  1,
		FirstName: "John",
		LastName:  "Doe",
		Dni:       "12345678",
		Gender:    "M",
		Phone:     "+5491123456789",
		Email:     "john.doe@example.com",
		BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	t.Run("✅ Success - Valid user update", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByDni", ctx, validUser.Dni).Return(nil, domain.ErrUserNotFound).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		// Act
		err := service.Update(ctx, validUser)

		// Assert
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("✅ Success - Update with same DNI", func(t *testing.T) {
		// Arrange
		existingUser := &domain.User{ID: 1, Dni: validUser.Dni}
		mockRepo.On("GetByDni", ctx, validUser.Dni).Return(existingUser, nil).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		// Act
		err := service.Update(ctx, validUser)

		// Assert
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("❌ Error - Nil user", func(t *testing.T) {
		// Act
		err := service.Update(ctx, nil)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user cannot be nil")
		mockRepo.AssertNotCalled(t, "GetByDni")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("❌ Error - Missing user ID", func(t *testing.T) {
		// Arrange
		userWithoutID := *validUser
		userWithoutID.ID = 0

		// Act
		err := service.Update(ctx, &userWithoutID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user id is required")
		mockRepo.AssertNotCalled(t, "GetByDni")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("❌ Error - DNI already exists for another user", func(t *testing.T) {
		// Arrange
		anotherUser := &domain.User{ID: 2, Dni: validUser.Dni}
		mockRepo.On("GetByDni", ctx, validUser.Dni).Return(anotherUser, nil).Once()

		// Act
		err := service.Update(ctx, validUser)

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrDniAlreadyExist))
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("❌ Error - Repository error on DNI check", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByDni", ctx, validUser.Dni).Return(nil, errors.New("database error")).Once()

		// Act
		err := service.Update(ctx, validUser)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error checking existing dni")
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("❌ Error - Repository error on update", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByDni", ctx, validUser.Dni).Return(nil, domain.ErrUserNotFound).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(errors.New("update error")).Once()

		// Act
		err := service.Update(ctx, validUser)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("❌ Error - Invalid user data - missing first name", func(t *testing.T) {
		// Arrange
		invalidUser := *validUser
		invalidUser.FirstName = ""

		// Act
		err := service.Update(ctx, &invalidUser)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "first name is required")
		mockRepo.AssertNotCalled(t, "GetByDni")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("❌ Error - Invalid user data - user is minor", func(t *testing.T) {
		// Arrange
		invalidUser := *validUser
		invalidUser.BirthDay = time.Now().AddDate(-17, 0, 0) // 17 years old

		// Act
		err := service.Update(ctx, &invalidUser)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user must be over 18")
		mockRepo.AssertNotCalled(t, "GetByDni")
		mockRepo.AssertNotCalled(t, "Update")
	})
}

func TestUserService_Delete(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepo()
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("✅ Success - Valid user deletion", func(t *testing.T) {
		// Arrange
		userID := uint(1)
		mockRepo.On("Delete", ctx, userID).Return(nil).Once()

		// Act
		err := service.Delete(ctx, userID)

		// Assert
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("❌ Error - Invalid user ID (zero)", func(t *testing.T) {
		// Act
		err := service.Delete(ctx, 0)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid user id")
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("❌ Error - User not found", func(t *testing.T) {
		// Arrange
		userID := uint(999)
		mockRepo.On("Delete", ctx, userID).Return(domain.ErrUserNotFound).Once()

		// Act
		err := service.Delete(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrUserNotFound))
		mockRepo.AssertExpectations(t)
	})

	t.Run("❌ Error - Repository error", func(t *testing.T) {
		// Arrange
		userID := uint(1)
		mockRepo.On("Delete", ctx, userID).Return(errors.New("database error")).Once()

		// Act
		err := service.Delete(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}
