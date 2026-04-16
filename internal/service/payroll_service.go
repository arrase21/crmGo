package service

import (
	"context"
	"errors"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
)

type PayrollService struct {
	payrollRepo domain.PayrollRepo
}

func NewPayrollService(u domain.PayrollRepo) *PayrollService {
	return &PayrollService{
		payrollRepo: u,
	}
}

func (s *PayrollService) Create(ctx context.Context, payroll *domain.Payroll) error {
	if payroll == nil {
		return errors.New("payroll cannot be nil")
	}
	return s.payrollRepo.Create(ctx, payroll)
}

func (s *PayrollService) GetByID(ctx context.Context, id uint) (*domain.Payroll, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	return s.payrollRepo.GetByID(ctx, id)
}

func (s *PayrollService) GetByEmployeeAndPeriod(ctx context.Context, employeID uint, periodStart, periodEnd time.Time) (*domain.Payroll, error) {
	if employeID == 0 {
		return nil, errors.New("invalid employee id")
	}
	if periodStart.IsZero() || periodEnd.IsZero() {
		return nil, errors.New("invalid period dates")
	}
	return s.payrollRepo.GetByEmployeeAndPeriod(ctx, employeID, periodStart, periodEnd)
}

func (s *PayrollService) ListByEmployee(ctx context.Context, employeeID uint) ([]domain.Payroll, error) {
	if employeeID == 0 {
		return nil, errors.New("invalid employee id")
	}
	return s.payrollRepo.ListByEmployee(ctx, employeeID)
}

func (s *PayrollService) Update(ctx context.Context, employee *domain.Payroll) error {
	if employee == nil {
		return errors.New("emplpyee cannot be nil")
	}
	if employee.ID == 0 {
		return errors.New("invalid employe id")
	}
	return s.payrollRepo.Update(ctx, employee)
}

func (s *PayrollService) Delete(ctx context.Context, employeeID uint) error {
	if employeeID == 0 {
		return errors.New("invalid employee id")
	}
	return s.payrollRepo.Delete(ctx, employeeID)
}

// GetByPeriod obtiene todas las nóminas de un periodo para el tenant actual
func (s *PayrollService) GetByPeriod(ctx context.Context, periodStart, periodEnd time.Time) ([]domain.Payroll, error) {
	if periodStart.IsZero() || periodEnd.IsZero() {
		return nil, errors.New("invalid period dates")
	}
	return s.payrollRepo.GetByPeriod(ctx, periodStart, periodEnd)
}
