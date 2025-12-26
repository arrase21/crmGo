package domain

import (
	"context"
	"errors"
)

// domain errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDniAlreadyExist   = errors.New("dni already exists")
	ErrEmailAlreadyExist = errors.New("email already exists") // ✅ Opcional: para futuro
	ErrPhoneAlreadyExist = errors.New("phone already exists") // ✅ Opcional: para futuro
	ErrInvalidTenantID   = errors.New("invalid tenant id")    // ✅ Opcional: para futuro
)

// ContextKey for tenant
type contextKey string

const TenantIDKey contextKey = "tenant_id"

type UserRepo interface {
	Create(ctx context.Context, usr *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByDni(ctx context.Context, dni string) (*User, error)
	List(ctx context.Context) ([]User, error)
	Update(ctx context.Context, usr *User) error
	Delete(ctx context.Context, id uint) error
}
