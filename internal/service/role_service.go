package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/arrase21/crm-users/internal/domain"
)

type RoleService struct {
	roleRepo domain.RoleRepo
}

func NewRoleService(r domain.RoleRepo) *RoleService {
	return &RoleService{
		roleRepo: r,
	}
}

func (s *RoleService) Create(ctx context.Context, role *domain.Role) error {
	if role == nil {
		return errors.New("role cannot be nil")
	}
	if role.Name == "" {
		return errors.New("role name is required")
	}

	existing, err := s.roleRepo.GetByName(ctx, role.Name)
	if err != nil && !errors.Is(err, domain.ErrRoleNotFound) {
		return fmt.Errorf("error checking existing role: %w", err)
	}
	if existing != nil {
		return domain.ErrRoleExisting
	}
	return s.roleRepo.Create(ctx, role)

}

func (s *RoleService) GetByID(ctx context.Context, id uint) (*domain.Role, error) {
	if id == 0 {
		return nil, errors.New("invalid role id")
	}
	return s.roleRepo.GetByID(ctx, id)
}

func (s *RoleService) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	if name == "" {
		return nil, errors.New("invalid name")
	}
	return s.roleRepo.GetByName(ctx, name)
}

func (s *RoleService) List(ctx context.Context) ([]domain.Role, error) {
	return s.roleRepo.List(ctx)
}

func (s *RoleService) Update(ctx context.Context, role *domain.Role) error {
	if role == nil || role.ID == 0 {
		return errors.New("invalid role")
	}
	_, err := s.roleRepo.GetByID(ctx, role.ID)
	if err != nil {
		return err
	}

	existingByname, err := s.roleRepo.GetByName(ctx, role.Name)
	if err != nil && !errors.Is(err, domain.ErrRoleNotFound) {
		return fmt.Errorf("error checking: %w", err)
	}
	if existingByname != nil && existingByname.ID != role.ID {
		return domain.ErrRoleExisting
	}
	return s.roleRepo.Update(ctx, role)
}

func (s *RoleService) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid role id")
	}
	return s.roleRepo.Delete(ctx, id)
}
