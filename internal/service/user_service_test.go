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

func TestUserService_Create(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepo()
	service := NewUserService(mockRepo)
	ctx := context.Background()

	// Test data
	validUser := &domain.User{
		TenantID:  1,
		FirstName: "John",
		LastName:  "Doe",
		Dni:       "12345678",
		Gender:    "M",
		Phone:     "+5491123456789",
		Email:     "john.doe@example.com",
		BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	t.Run("✅ Success - Valid user creation", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByDni", ctx, validUser.Dni).Return(nil, domain.ErrUserNotFound).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		// Act
		err := service.Create(ctx, validUser)

		// Assert
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("❌ Error - Nil user", func(t *testing.T) {
		// Act
		err := service.Create(ctx, nil)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user cannot be nil")
		mockRepo.AssertNotCalled(t, "GetByDni")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("❌ Error - DNI already exists", func(t *testing.T) {
		// Arrange
		existingUser := &domain.User{ID: 1, Dni: validUser.Dni}
		mockRepo.On("GetByDni", ctx, validUser.Dni).Return(existingUser, nil).Once()

		// Act
		err := service.Create(ctx, validUser)

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrDniAlreadyExist))
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("❌ Error - Repository error on DNI check", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByDni", ctx, validUser.Dni).Return(nil, errors.New("database error")).Once()

		// Act
		err := service.Create(ctx, validUser)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error checking existing user")
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("❌ Error - Repository error on create", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByDni", ctx, validUser.Dni).Return(nil, domain.ErrUserNotFound).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(errors.New("create error")).Once()

		// Act
		err := service.Create(ctx, validUser)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("❌ Error - Invalid user data - missing first name", func(t *testing.T) {
		// Arrange
		invalidUser := *validUser
		invalidUser.FirstName = ""

		// Act
		err := service.Create(ctx, &invalidUser)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "first name is required")
		mockRepo.AssertNotCalled(t, "GetByDni")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("❌ Error - Invalid user data - invalid gender", func(t *testing.T) {
		// Arrange
		invalidUser := *validUser
		invalidUser.Gender = "X"

		// Act
		err := service.Create(ctx, &invalidUser)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid option")
		mockRepo.AssertNotCalled(t, "GetByDni")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("❌ Error - Invalid user data - user is minor", func(t *testing.T) {
		// Arrange
		invalidUser := *validUser
		invalidUser.BirthDay = time.Now().AddDate(-17, 0, 0) // 17 years old

		// Act
		err := service.Create(ctx, &invalidUser)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user must be over 18")
		mockRepo.AssertNotCalled(t, "GetByDni")
		mockRepo.AssertNotCalled(t, "Create")
	})
}

func TestUserService_GetByID(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepo()
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("✅ Success - Valid user ID", func(t *testing.T) {
		// Arrange
		expectedUser := &domain.User{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
		}
		mockRepo.On("GetByID", ctx, uint(1)).Return(expectedUser, nil).Once()

		// Act
		user, err := service.GetByID(ctx, 1)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("❌ Error - Invalid user ID (zero)", func(t *testing.T) {
		// Act
		user, err := service.GetByID(ctx, 0)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid user id")
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("❌ Error - User not found", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByID", ctx, uint(999)).Return(nil, domain.ErrUserNotFound).Once()

		// Act
		user, err := service.GetByID(ctx, 999)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, errors.Is(err, domain.ErrUserNotFound))
		mockRepo.AssertExpectations(t)
	})

	t.Run("❌ Error - Repository error", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("database error")).Once()

		// Act
		user, err := service.GetByID(ctx, 1)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByDni(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepo()
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("✅ Success - Valid DNI", func(t *testing.T) {
		// Arrange
		dni := "12345678"
		expectedUser := &domain.User{
			ID:        1,
			Dni:       dni,
			FirstName: "John",
			LastName:  "Doe",
		}
		mockRepo.On("GetByDni", ctx, dni).Return(expectedUser, nil).Once()

		// Act
		user, err := service.GetByDni(ctx, dni)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("❌ Error - Empty DNI", func(t *testing.T) {
		// Act
		user, err := service.GetByDni(ctx, "")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "dni cannot be empty")
		mockRepo.AssertNotCalled(t, "GetByDni")
	})

	t.Run("❌ Error - User not found", func(t *testing.T) {
		// Arrange
		dni := "99999999"
		mockRepo.On("GetByDni", ctx, dni).Return(nil, domain.ErrUserNotFound).Once()

		// Act
		user, err := service.GetByDni(ctx, dni)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, errors.Is(err, domain.ErrUserNotFound))
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_List(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepo()
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("✅ Success - List users", func(t *testing.T) {
		// Arrange
		expectedUsers := []domain.User{
			{ID: 1, FirstName: "John", LastName: "Doe"},
			{ID: 2, FirstName: "Jane", LastName: "Smith"},
		}
		mockRepo.On("List", ctx, 1, 20).Return(expectedUsers, int64(2), nil).Once()

		// Act
		users, total, err := service.List(ctx, 1, 20)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Equal(t, int64(2), total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("❌ Error - Repository error", func(t *testing.T) {
		// Arrange
		mockRepo.On("List", ctx, 1, 20).Return([]domain.User(nil), int64(0), errors.New("database error")).Once()

		// Act
		users, total, err := service.List(ctx, 1, 20)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Equal(t, int64(0), total)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}
