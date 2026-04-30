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
	ErrAbsenceNotFound = errors.New("absence not found")
	ErrAbsenceInvalid  = errors.New("invalid absence data")
)

// ========================================
// Absence Service
// ========================================

type AbsenceService struct {
	absenceRepo  domain.AbsenceRepo
	employeeRepo domain.EmployeeRepo
}

func NewAbsenceService(absenceRepo domain.AbsenceRepo, employeeRepo domain.EmployeeRepo) *AbsenceService {
	return &AbsenceService{
		absenceRepo:  absenceRepo,
		employeeRepo: employeeRepo,
	}
}

// AbsenceRequest request para crear una ausencia
type AbsenceRequest struct {
	EmployeeID         uint    `json:"employee_id" validate:"required"`
	Type               string  `json:"type" validate:"required,oneof=sick maternity paternity vacation unpaid other"`
	StartDate          string  `json:"start_date" validate:"required"`
	EndDate            string  `json:"end_date" validate:"required"`
	PaidPercent        float64 `json:"paid_percent"` // default 100
	MedicalCertificate string  `json:"medical_certificate"`
	Notes              string  `json:"notes"`
}

// CalculateAbsenceDeduction calcula la deducción por ausencia
func (s *AbsenceService) CalculateAbsenceDeduction(baseSalary float64, absence *domain.Absence) float64 {
	days := int(absence.EndDate.Sub(absence.StartDate).Hours()/24) + 1
	dailyRate := baseSalary / 30

	// Porcentaje no pagado
	unpaidPercent := 100 - absence.PaidPercent
	return dailyRate * float64(days) * (unpaidPercent / 100)
}

// Create crea un registro de ausencia
func (s *AbsenceService) Create(ctx context.Context, req AbsenceRequest) (*domain.Absence, error) {
	// Validar empleado existe
	emp, err := s.employeeRepo.GetByID(ctx, req.EmployeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Parsear fechas
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errors.New("invalid start_date format, use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, errors.New("invalid end_date format, use YYYY-MM-DD")
	}

	if endDate.Before(startDate) {
		return nil, errors.New("end_date cannot be before start_date")
	}

	paidPercent := req.PaidPercent
	if paidPercent == 0 {
		// Por defecto, las vacaciones son 100%, sick puede variar
		switch req.Type {
		case domain.AbsenceTypeVacation:
			paidPercent = 100
		case domain.AbsenceTypeSick:
			paidPercent = 66.6 // Ejemplo: 2/3 del salario
		default:
			paidPercent = 0
		}
	}

	// Obtener contrato para calcular deducción
	contract, _ := s.getActiveContract(ctx, req.EmployeeID)
	var dailyRate, totalDeduction float64
	if contract != nil {
		dailyRate = contract.BaseSalary / 30
		days := int(endDate.Sub(startDate).Hours()/24) + 1
		unpaidPercent := 100 - paidPercent
		totalDeduction = dailyRate * float64(days) * (unpaidPercent / 100)
	}

	absence := &domain.Absence{
		EmployeeID:         req.EmployeeID,
		Type:               req.Type,
		StartDate:          startDate,
		EndDate:            endDate,
		PaidPercent:        paidPercent,
		DailyRate:          dailyRate,
		TotalDeduction:     totalDeduction,
		MedicalCertificate: req.MedicalCertificate,
		Status:             "active",
		Notes:              req.Notes,
	}

	if err := s.absenceRepo.Create(ctx, absence); err != nil {
		return nil, err
	}

	absence.Employee = *emp
	return absence, nil
}

// Process marca una ausencia como procesada
func (s *AbsenceService) Process(ctx context.Context, id uint, payrollID uint) (*domain.Absence, error) {
	absence, err := s.absenceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrAbsenceNotFound
	}

	if absence.Status == "processed" {
		return nil, errors.New("absence already processed")
	}

	absence.Status = "processed"
	absence.PayrollID = payrollID

	if err := s.absenceRepo.Update(ctx, absence); err != nil {
		return nil, err
	}

	return absence, nil
}

// ListByEmployee lista ausencias de un empleado
func (s *AbsenceService) ListByEmployee(ctx context.Context, employeeID uint, page, limit int) ([]domain.Absence, int64, error) {
	return s.absenceRepo.List(ctx, employeeID, page, limit)
}

// GetByPayroll obtiene las ausencias de una nómina
func (s *AbsenceService) GetByPayroll(ctx context.Context, payrollID uint) ([]domain.Absence, error) {
	return s.absenceRepo.GetByPayrollID(ctx, payrollID)
}

// GetTotalDeductionByPayroll obtiene el total de deducciones por ausencias
func (s *AbsenceService) GetTotalDeductionByPayroll(ctx context.Context, payrollID uint) float64 {
	absences, _ := s.absenceRepo.GetByPayrollID(ctx, payrollID)
	var total float64
	for _, a := range absences {
		total += a.TotalDeduction
	}
	return total
}

func (s *AbsenceService) getActiveContract(ctx context.Context, employeeID uint) (*domain.EmployeeContract, error) {
	emp, err := s.employeeRepo.GetByID(ctx, employeeID)
	if err != nil {
		return nil, err
	}
	return emp.GetMainContract(), nil
}
