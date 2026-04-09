package service

import (
	"context"
	"errors"

	"github.com/arrase21/crm-users/internal/domain"
)

type EmployeeService struct {
	empRepo domain.EmployeeRepo
}

func NewEmployeeService(u domain.EmployeeRepo) *EmployeeService {
	return &EmployeeService{
		empRepo: u,
	}
}

func (s *EmployeeService) Create(ctx context.Context, emp *domain.Employee) error {
	if emp == nil {
		return errors.New("employee cannot be nil")
	}
	return s.empRepo.Create(ctx, emp)
}

func (s *EmployeeService) GetByID(ctx context.Context, id uint) (*domain.Employee, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	return s.empRepo.GetByID(ctx, id)
}

func (s *EmployeeService) GetByUserID(ctx context.Context, userID uint) (*domain.Employee, error) {
	if userID == 0 {
		return nil, errors.New("Invalid userID")
	}
	return s.empRepo.GetByUserID(ctx, userID)
}

func (s *EmployeeService) List(ctx context.Context, page, limit int) ([]domain.Employee, int64, error) {
	return s.empRepo.List(ctx, page, limit)
}

func (s *EmployeeService) Update(ctx context.Context, emp *domain.Employee) error {
	if emp == nil {
		return errors.New("employee cannot be nil")
	}
	if emp.ID == 0 {
		return errors.New("invalid employee id")
	}
	return s.empRepo.Update(ctx, emp)
}

func (s *EmployeeService) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid id")
	}
	return s.empRepo.Delete(ctx, id)
}
