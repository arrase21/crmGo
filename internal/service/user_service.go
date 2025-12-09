package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/arrase21/crm-users/internal/domain"
)

type UserService struct {
	usrRepo domain.UserRepo
}

func NewUserService(u domain.UserRepo) *UserService {
	return &UserService{
		usrRepo: u,
	}
}

func (s *UserService) Create(ctx context.Context, usr *domain.User) error {
	// validate
	usr.Normalize()
	if err := usr.Validate(); err != nil {
		return fmt.Errorf("Validation error in domain %w", err)
	}
	// context
	existing, err := s.usrRepo.GetByDni(ctx, usr.Dni)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return fmt.Errorf("error checking existing user: %w", err)
	}
	if existing != nil {
		return errors.New("users with this dni already exixting")
	}
	return s.usrRepo.Create(ctx, usr)
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user id")
	}
	return s.usrRepo.GetByID(ctx, id)
}

func (s *UserService) GetByDni(ctx context.Context, dni string) (*domain.User, error) {
	return s.usrRepo.GetByDni(ctx, dni)
}

func (s *UserService) List(ctx context.Context) ([]domain.User, error) {
	return s.usrRepo.List(ctx)
}

func (s *UserService) Update(ctx context.Context, usr *domain.User) error {
	if usr.ID == 0 {
		return errors.New("user id is required")
	}
	// validate
	usr.Normalize()
	if err := usr.Validate(); err != nil {
		return fmt.Errorf("validation erro in domain: %w", err)
	}
	// duplicates
	existing, err := s.usrRepo.GetByDni(ctx, usr.Dni)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return err
	}
	if existing != nil && existing.ID != usr.ID {
		return errors.New("Dni already exists for another user")
	}
	return s.usrRepo.Update(ctx, usr)
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("Invalid user id")
	}
	return s.usrRepo.Delete(ctx, id)
}
