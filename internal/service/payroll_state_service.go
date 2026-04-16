package service

import (
	"context"
	"errors"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
)

// Errores de transición de estado
var (
	ErrInvalidStatusTransition = errors.New("invalid status transition: cannot change from current status")
	ErrPayrollNotInDraft       = errors.New("payroll must be in draft status to modify")
)

// PayrollStateService maneja las transiciones de estado de la nómina
type PayrollStateService struct {
	payrollRepo  domain.PayrollRepo
	paymentRepo  domain.PaymentRepo
	employeeRepo domain.EmployeeRepo
}

func NewPayrollStateService(
	payrollRepo domain.PayrollRepo,
	paymentRepo domain.PaymentRepo,
	employeeRepo domain.EmployeeRepo,
) *PayrollStateService {
	return &PayrollStateService{
		payrollRepo:  payrollRepo,
		paymentRepo:  paymentRepo,
		employeeRepo: employeeRepo,
	}
}

// MarkAsPaid marca una nómina como pagada y crea el registro de pago
func (s *PayrollStateService) MarkAsPaid(ctx context.Context, payrollID uint, paymentMethod string) (*domain.Payment, error) {
	if payrollID == 0 {
		return nil, errors.New("payroll id is required")
	}

	// Obtener la nómina
	payroll, err := s.payrollRepo.GetByID(ctx, payrollID)
	if err != nil {
		return nil, err
	}

	// Verificar que no esté ya pagada (primero, mensaje más específico)
	if payroll.Status == domain.PayrollStatusPaid {
		return nil, errors.New("payroll is already paid")
	}

	// Validar que esté en estado válido para transicionar
	if !s.canTransitionTo(payroll.Status, domain.PayrollStatusPaid) {
		return nil, ErrInvalidStatusTransition
	}

	// Crear el registro de pago
	payment := &domain.Payment{
		PayrollID: payroll.ID,
		Method:    paymentMethod,
		BankName:  "", // Se puede obtener del contrato si existe
		Amount:    payroll.NetAmount,
		PaidAt:    time.Now(),
		Status:    "completed",
		CreatedAt: time.Now(),
	}

	// Guardar el pago
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	// Actualizar estado de la nómina
	payroll.Status = domain.PayrollStatusPaid
	if err := s.payrollRepo.Update(ctx, payroll); err != nil {
		// Intentar rollback del pago
		s.paymentRepo.Delete(ctx, payment.ID)
		return nil, err
	}

	return payment, nil
}

// MarkAsCalculated marca una nómina draft como calculada
func (s *PayrollStateService) MarkAsCalculated(ctx context.Context, payrollID uint) error {
	if payrollID == 0 {
		return errors.New("payroll id is required")
	}

	payroll, err := s.payrollRepo.GetByID(ctx, payrollID)
	if err != nil {
		return err
	}

	if !s.canTransitionTo(payroll.Status, domain.PayrollStatusCalculated) {
		return ErrInvalidStatusTransition
	}

	payroll.Status = domain.PayrollStatusCalculated
	return s.payrollRepo.Update(ctx, payroll)
}

// RevertToDraft revierte una nómina calculada (no pagada) a draft
func (s *PayrollStateService) RevertToDraft(ctx context.Context, payrollID uint) error {
	if payrollID == 0 {
		return errors.New("payroll id is required")
	}

	payroll, err := s.payrollRepo.GetByID(ctx, payrollID)
	if err != nil {
		return err
	}

	// Solo se puede revertir si está calculada (no pagada)
	if payroll.Status == domain.PayrollStatusPaid {
		return errors.New("cannot revert a paid payroll")
	}

	if !s.canTransitionTo(payroll.Status, domain.PayrollStatusDraft) {
		return ErrInvalidStatusTransition
	}

	payroll.Status = domain.PayrollStatusDraft
	return s.payrollRepo.Update(ctx, payroll)
}

// ValidatePayrollCreation valida que no exista una nómina para el mismo periodo
func (s *PayrollStateService) ValidatePayrollCreation(ctx context.Context, employeeID uint, periodStart, periodEnd time.Time) error {
	existing, err := s.payrollRepo.GetByEmployeeAndPeriod(ctx, employeeID, periodStart, periodEnd)
	if err != nil {
		if errors.Is(err, domain.ErrPayrollNotFound) {
			return nil // No existe, se puede crear
		}
		return err
	}

	// Si existe y está pagada, no se puede recalcular
	if existing.Status == domain.PayrollStatusPaid {
		return errors.New("payroll already paid for this period, cannot recalculate")
	}

	// Si existe y está calculada, se permite reescribir
	return nil
}

// canTransitionTo valida si una transición de estado es válida
// Estados posibles: draft -> calculated -> paid
// Excepciones: calculated -> draft (revertir)
func (s *PayrollStateService) canTransitionTo(from, to string) bool {
	validTransitions := map[string][]string{
		domain.PayrollStatusDraft:      {domain.PayrollStatusCalculated, domain.PayrollStatusDraft},
		domain.PayrollStatusCalculated: {domain.PayrollStatusPaid, domain.PayrollStatusDraft},
		domain.PayrollStatusPaid:       {}, // Estado terminal, no hay transiciones válidas
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, state := range allowed {
		if state == to {
			return true
		}
	}
	return false
}

// GetPaymentHistory obtiene el historial de pagos de una nómina
func (s *PayrollStateService) GetPaymentHistory(ctx context.Context, payrollID uint) (*domain.Payment, error) {
	return s.paymentRepo.GetByPayrollID(ctx, payrollID)
}
