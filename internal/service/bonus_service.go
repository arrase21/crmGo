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
	ErrBonusNotFound = errors.New("bonus not found")
	ErrBonusInvalid  = errors.New("invalid bonus data")
)

// ========================================
// Bonus Service
// ========================================

type BonusService struct {
	bonusRepo    domain.BonusRepo
	employeeRepo domain.EmployeeRepo
}

func NewBonusService(bonusRepo domain.BonusRepo, employeeRepo domain.EmployeeRepo) *BonusService {
	return &BonusService{
		bonusRepo:    bonusRepo,
		employeeRepo: employeeRepo,
	}
}

// BonusRequest request para crear una bonificación
type BonusRequest struct {
	EmployeeID  uint    `json:"employee_id" validate:"required"`
	Type        string  `json:"type" validate:"required,oneof=performance production attendance christmas other"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description string  `json:"description"`
	Date        string  `json:"date" validate:"required"`
	Notes       string  `json:"notes"`
}

// Create crea un registro de bonificación
func (s *BonusService) Create(ctx context.Context, req BonusRequest) (*domain.Bonus, error) {
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

	bonus := &domain.Bonus{
		EmployeeID:  req.EmployeeID,
		Type:        req.Type,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        date,
		Status:      "pending",
		Notes:       req.Notes,
	}

	if err := s.bonusRepo.Create(ctx, bonus); err != nil {
		return nil, err
	}

	bonus.Employee = *emp
	return bonus, nil
}

// Approve marca una bonificación como aprobada
func (s *BonusService) Approve(ctx context.Context, id uint, approvedBy uint) (*domain.Bonus, error) {
	bonus, err := s.bonusRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrBonusNotFound
	}

	if bonus.Status == "approved" || bonus.Status == "paid" {
		return nil, errors.New("bonus already approved or paid")
	}

	bonus.Status = "approved"
	bonus.ApprovedBy = approvedBy
	now := time.Now()
	bonus.ApprovedAt = &now

	if err := s.bonusRepo.Update(ctx, bonus); err != nil {
		return nil, err
	}

	return bonus, nil
}

// ListByEmployee lista bonificaciones de un empleado
func (s *BonusService) ListByEmployee(ctx context.Context, employeeID uint, page, limit int) ([]domain.Bonus, int64, error) {
	return s.bonusRepo.List(ctx, employeeID, page, limit)
}

// GetByPayroll obtiene las bonificaciones de una nómina
func (s *BonusService) GetByPayroll(ctx context.Context, payrollID uint) ([]domain.Bonus, error) {
	return s.bonusRepo.GetByPayrollID(ctx, payrollID)
}

// GetTotalAmountByPayroll obtiene el total de bonificaciones de una nómina
func (s *BonusService) GetTotalAmountByPayroll(ctx context.Context, payrollID uint) float64 {
	bonuses, _ := s.bonusRepo.GetByPayrollID(ctx, payrollID)
	var total float64
	for _, b := range bonuses {
		total += b.Amount
	}
	return total
}
