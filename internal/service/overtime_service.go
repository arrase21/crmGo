package service

import (
	"context"
	"errors"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
)

// ========================================
// Errors
// ========================================

var (
	ErrOvertimeNotFound    = errors.New("overtime not found")
	ErrOvertimeInvalid     = errors.New("invalid overtime data")
	ErrOvertimeNotApproved = errors.New("overtime not approved")
)

// ========================================
// Overtime Service
// ========================================

type OvertimeService struct {
	overtimeRepo domain.OvertimeRepo
	employeeRepo domain.EmployeeRepo
}

func NewOvertimeService(overtimeRepo domain.OvertimeRepo, employeeRepo domain.EmployeeRepo) *OvertimeService {
	return &OvertimeService{
		overtimeRepo: overtimeRepo,
		employeeRepo: employeeRepo,
	}
}

// OvertimeRequest request para crear hora extra
type OvertimeRequest struct {
	EmployeeID uint    `json:"employee_id" validate:"required"`
	Date       string  `json:"date" validate:"required"`
	Hours      float64 `json:"hours" validate:"required,gt=0"`
	Type       string  `json:"type" validate:"required,oneof=extra night holiday sunday extra_night"`
	Notes      string  `json:"notes"`
}

// CalculateOvertimeAmount calcula el valor de las horas extras
func (s *OvertimeService) CalculateOvertimeAmount(baseSalary float64, hours float64, otType string) float64 {
	hourlyRate := baseSalary / 240 // Assuming 240 working hours/month

	rate := 1.0
	switch otType {
	case domain.OvertimeTypeExtra:
		rate = 1.25 // 25% recargo
	case domain.OvertimeTypeNight:
		rate = 1.35 // 35% recargo nocturno
	case domain.OvertimeTypeHoliday:
		rate = 2.0 // 100% festivo
	case domain.OvertimeTypeSunday:
		rate = 1.75 // 75% dominical
	case domain.OvertimeTypeExtraNight:
		rate = 2.25 // 125% extra nocturna
	}

	return hourlyRate * hours * rate
}

// Create crea un registro de hora extra
func (s *OvertimeService) Create(ctx context.Context, req OvertimeRequest) (*domain.Overtime, error) {
	// Validar empleado existe
	emp, err := s.employeeRepo.GetByID(ctx, req.EmployeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Parsear fecha
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	// Obtener contrato para calcular monto
	contract, err := s.getActiveContract(ctx, req.EmployeeID)
	var amount float64
	if err == nil && contract != nil {
		amount = s.CalculateOvertimeAmount(contract.BaseSalary, req.Hours, req.Type)
	}

	overtime := &domain.Overtime{
		EmployeeID: req.EmployeeID,
		Date:       date,
		Hours:      req.Hours,
		Type:       req.Type,
		Amount:     amount,
		Status:     "pending",
		Notes:      req.Notes,
	}

	if err := s.overtimeRepo.Create(ctx, overtime); err != nil {
		return nil, err
	}

	overtime.Employee = *emp
	return overtime, nil
}

// Approve marca una hora extra como aprobada
func (s *OvertimeService) Approve(ctx context.Context, id uint, approvedBy uint) (*domain.Overtime, error) {
	overtime, err := s.overtimeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrOvertimeNotFound
	}

	if overtime.Status == "approved" {
		return nil, errors.New("overtime already approved")
	}

	overtime.Status = "approved"
	overtime.ApprovedBy = approvedBy
	now := time.Now()
	overtime.ApprovedAt = &now

	if err := s.overtimeRepo.Update(ctx, overtime); err != nil {
		return nil, err
	}

	return overtime, nil
}

// ListByEmployee lista horas extras de un empleado
func (s *OvertimeService) ListByEmployee(ctx context.Context, employeeID uint, page, limit int) ([]domain.Overtime, int64, error) {
	return s.overtimeRepo.List(ctx, employeeID, page, limit)
}

// GetByPayroll obtiene las horas extras de una nómina
func (s *OvertimeService) GetByPayroll(ctx context.Context, payrollID uint) ([]domain.Overtime, error) {
	return s.overtimeRepo.GetByPayrollID(ctx, payrollID)
}

// GetTotalAmountByPayroll obtiene el total de horas extras de una nómina
func (s *OvertimeService) GetTotalAmountByPayroll(ctx context.Context, payrollID uint) float64 {
	overtimes, _ := s.overtimeRepo.GetByPayrollID(ctx, payrollID)
	var total float64
	for _, ot := range overtimes {
		total += ot.Amount
	}
	return total
}

func (s *OvertimeService) getActiveContract(ctx context.Context, employeeID uint) (*domain.EmployeeContract, error) {
	emp, err := s.employeeRepo.GetByID(ctx, employeeID)
	if err != nil {
		return nil, err
	}
	return emp.GetMainContract(), nil
}
