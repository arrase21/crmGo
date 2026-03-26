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
		return nil, errors.New("invalid user id")
	}
	return s.empRepo.GetByID(ctx, id)
}
