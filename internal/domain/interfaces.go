package domain

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepo interface {
	Create(ctx context.Context, usr *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByDni(ctx context.Context, dni string) (*User, error)
	List(ctx context.Context) ([]User, error)
	Update(ctx context.Context, usr *User) error
	Delete(ctx context.Context, id uint) error
}
