package service

import (
	"context"
	"errors"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
)

type PayrollCalculatorService struct {
	payrollRepo        domain.PayrollRepo
	payrollItemRepo    domain.PayrollItemRepo
	employeeRepo       domain.EmployeeRepo
	contractRepo       domain.EmployeeContractRepo
	payrollConceptRepo domain.PayrollConceptRepo
}

func NewPayrollCalculatorService(
	payrollRepo domain.PayrollRepo,
	payrollItemRepo domain.PayrollItemRepo,
	employeeRepo domain.EmployeeRepo,
	contractRepo domain.EmployeeContractRepo,
	conceptRepo domain.PayrollConceptRepo,
) *PayrollCalculatorService {
	return &PayrollCalculatorService{
		payrollRepo:        payrollRepo,
		payrollItemRepo:    payrollItemRepo,
		employeeRepo:       employeeRepo,
		contractRepo:       contractRepo,
		payrollConceptRepo: conceptRepo,
	}
}

type CalculatePayrollRequest struct {
	EmployeeID  uint
	PeriodStart time.Time
	PeriodEnd   time.Time
	PayDate     time.Time
}

type CalculatedPayroll struct {
	Payroll         *domain.Payroll
	Items           []domain.PayrollItem
	GrossAmount     float64
	TotalDeductions float64
	NetAmount       float64
}

func (s *PayrollCalculatorService) Calculate(ctx context.Context, req CalculatePayrollRequest) (*CalculatedPayroll, error) {
	if req.EmployeeID == 0 {
		return nil, errors.New("employee id is required")
	}
	if req.PeriodStart.IsZero() || req.PeriodEnd.IsZero() {
		return nil, errors.New("period start and end are required")
	}
	if req.PeriodEnd.Before(req.PeriodStart) {
		return nil, errors.New("period end cannot be before period start")
	}

	employee, err := s.employeeRepo.GetByID(ctx, req.EmployeeID)
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			return nil, errors.New("employee not found")
		}
		return nil, err
	}

	contract, err := s.contractRepo.GetActiveByEmployee(ctx, req.EmployeeID)
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeContractNotFound) {
			return nil, errors.New("no active contract found for employee")
		}
		return nil, err
	}

	concepts, err := s.payrollConceptRepo.GetActiveConcepts(ctx)
	if err != nil {
		return nil, err
	}
	if len(concepts) == 0 {
		return nil, errors.New("no active payroll concepts configured")
	}

	baseSalary := contract.BaseSalary
	periodDays := int(req.PeriodEnd.Sub(req.PeriodStart).Hours()/24) + 1
	monthDays := 30.0
	if periodDays != 30 {
		baseSalary = baseSalary * float64(periodDays) / monthDays
	}

	var items []domain.PayrollItem
	var grossAmount float64
	var totalDeductions float64

	for _, concept := range concepts {
		item := s.calculateConceptItem(concept, baseSalary, contract)
		items = append(items, item)

		switch concept.Type {
		case domain.PayrollTypeEarning:
			grossAmount += item.Amount
		case domain.PayrollTypeDeduction:
			totalDeductions += item.Amount
		case domain.PayrollTypeEmployerContribution:
		}
	}

	netAmount := grossAmount - totalDeductions

	payroll := &domain.Payroll{
		TenantID:        employee.TenantID,
		EmployeeID:      employee.ID,
		PeriodStart:     req.PeriodStart,
		PeriodEnd:       req.PeriodEnd,
		PayDate:         req.PayDate,
		GrossAmount:     grossAmount,
		TotalDeductions: totalDeductions,
		NetAmount:       netAmount,
		Status:          domain.PayrollStatusDraft,
	}

	return &CalculatedPayroll{
		Payroll:         payroll,
		Items:           items,
		GrossAmount:     grossAmount,
		TotalDeductions: totalDeductions,
		NetAmount:       netAmount,
	}, nil
}

func (s *PayrollCalculatorService) calculateConceptItem(
	concept domain.PayrollConcept,
	baseSalary float64,
	contract *domain.EmployeeContract,
) domain.PayrollItem {
	var amount float64

	switch concept.Type {
	case domain.PayrollTypeEarning:
		amount = s.calculateEarning(concept, baseSalary, contract)

	case domain.PayrollTypeDeduction:
		amount = baseSalary * concept.Percentage / 100

	case domain.PayrollTypeEmployerContribution:
		amount = baseSalary * concept.Percentage / 100
	}

	return domain.PayrollItem{
		ConceptID:    concept.ID,
		Type:         concept.Type,
		Code:         concept.Code,
		Name:         concept.Name,
		Amount:       amount,
		CalculatedAt: time.Now(),
	}
}

func (s *PayrollCalculatorService) calculateEarning(
	concept domain.PayrollConcept,
	baseSalary float64,
	contract *domain.EmployeeContract,
) float64 {
	if concept.Percentage > 0 {
		return baseSalary * concept.Percentage / 100
	}

	if concept.EmployeePart > 0 {
		return concept.EmployeePart
	}

	switch concept.Code {
	case domain.ConceptTransport:
		return contract.TransportAllowance
	case domain.ConceptHousing:
		return contract.HousingAllowance
	}
	return concept.EmployeePart
}

func (s *PayrollCalculatorService) CalculateAndSave(ctx context.Context, req CalculatePayrollRequest) (*CalculatedPayroll, error) {
	calculated, err := s.Calculate(ctx, req)
	if err != nil {
		return nil, err
	}

	existing, err := s.payrollRepo.GetByEmployeeAndPeriod(ctx, req.EmployeeID, req.PeriodStart, req.PeriodEnd)
	if err == nil {
		calculated.Payroll.ID = existing.ID
		err = s.payrollRepo.Update(ctx, calculated.Payroll)
		if err != nil {
			return nil, err
		}
		s.payrollItemRepo.DeleteByPayrollID(ctx, existing.ID)
	} else {
		err = s.payrollRepo.Create(ctx, calculated.Payroll)
		if err != nil {
			return nil, err
		}
	}
	var itemsToSave []domain.PayrollItem
	for i := range calculated.Items {
		calculated.Items[i].PayrollID = calculated.Payroll.ID
		itemsToSave = append(itemsToSave, calculated.Items[i])
	}
	err = s.payrollItemRepo.CreateBatch(ctx, itemsToSave)
	if err != nil {
		return nil, err
	}
	return calculated, nil
}

func (s *PayrollCalculatorService) CalculatePeriodSummary(ctx context.Context, periodStart, periodEnd time.Time) ([]CalculatedPayroll, error) {
	return nil, errors.New("not implemented: need list all employees endpoint")
}
