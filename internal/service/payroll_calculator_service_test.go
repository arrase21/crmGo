package service

import (
	"context"
	"testing"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ========================================
// Mocks para testing
// ========================================

type MockPayrollRepo struct {
	mock.Mock
}

func (m *MockPayrollRepo) Create(ctx context.Context, payroll *domain.Payroll) error {
	args := m.Called(ctx, payroll)
	if payroll != nil {
		payroll.ID = 1 // Simular ID asignado
	}
	return args.Error(0)
}

func (m *MockPayrollRepo) GetByID(ctx context.Context, id uint) (*domain.Payroll, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payroll), args.Error(1)
}

func (m *MockPayrollRepo) GetByEmployeeAndPeriod(ctx context.Context, employeeID uint, periodStart, periodEnd time.Time) (*domain.Payroll, error) {
	args := m.Called(ctx, employeeID, periodStart, periodEnd)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payroll), args.Error(1)
}

func (m *MockPayrollRepo) GetByPeriod(ctx context.Context, periodStart, periodEnd time.Time) ([]domain.Payroll, error) {
	args := m.Called(ctx, periodStart, periodEnd)
	return args.Get(0).([]domain.Payroll), args.Error(1)
}

func (m *MockPayrollRepo) ListByEmployee(ctx context.Context, employeeID uint) ([]domain.Payroll, error) {
	args := m.Called(ctx, employeeID)
	return args.Get(0).([]domain.Payroll), args.Error(1)
}

func (m *MockPayrollRepo) Update(ctx context.Context, payroll *domain.Payroll) error {
	args := m.Called(ctx, payroll)
	return args.Error(0)
}

func (m *MockPayrollRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockPayrollItemRepo struct {
	mock.Mock
}

func (m *MockPayrollItemRepo) Create(ctx context.Context, item *domain.PayrollItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockPayrollItemRepo) CreateBatch(ctx context.Context, items []domain.PayrollItem) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *MockPayrollItemRepo) GetByIDPayrollID(ctx context.Context, payrollID uint) ([]domain.PayrollItem, error) {
	args := m.Called(ctx, payrollID)
	return args.Get(0).([]domain.PayrollItem), args.Error(1)
}

func (m *MockPayrollItemRepo) DeleteByPayrollID(ctx context.Context, payrollID uint) error {
	args := m.Called(ctx, payrollID)
	return args.Error(0)
}

type MockEmployeeRepo struct {
	mock.Mock
}

func (m *MockEmployeeRepo) Create(ctx context.Context, emp *domain.Employee) error {
	args := m.Called(ctx, emp)
	return args.Error(0)
}

func (m *MockEmployeeRepo) GetByID(ctx context.Context, id uint) (*domain.Employee, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *MockEmployeeRepo) GetByUserID(ctx context.Context, userID uint) (*domain.Employee, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *MockEmployeeRepo) List(ctx context.Context, page, limit int) ([]domain.Employee, int64, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]domain.Employee), args.Get(1).(int64), args.Error(2)
}

func (m *MockEmployeeRepo) ListActive(ctx context.Context, page, limit int) ([]domain.Employee, int64, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]domain.Employee), args.Get(1).(int64), args.Error(2)
}

func (m *MockEmployeeRepo) Update(ctx context.Context, emp *domain.Employee) error {
	args := m.Called(ctx, emp)
	return args.Error(0)
}

func (m *MockEmployeeRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockContractRepo struct {
	mock.Mock
}

func (m *MockContractRepo) Create(ctx context.Context, contract *domain.EmployeeContract) error {
	args := m.Called(ctx, contract)
	return args.Error(0)
}

func (m *MockContractRepo) GetByID(ctx context.Context, id uint) (*domain.EmployeeContract, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EmployeeContract), args.Error(1)
}

func (m *MockContractRepo) GetActiveByEmployee(ctx context.Context, employeeID uint) (*domain.EmployeeContract, error) {
	args := m.Called(ctx, employeeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EmployeeContract), args.Error(1)
}

func (m *MockContractRepo) ListByEmployee(ctx context.Context, employeeID uint) ([]domain.EmployeeContract, error) {
	args := m.Called(ctx, employeeID)
	return args.Get(0).([]domain.EmployeeContract), args.Error(1)
}

func (m *MockContractRepo) Update(ctx context.Context, contract *domain.EmployeeContract) error {
	args := m.Called(ctx, contract)
	return args.Error(0)
}

func (m *MockContractRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockConceptRepo struct {
	mock.Mock
}

func (m *MockConceptRepo) Create(ctx context.Context, concept *domain.PayrollConcept) error {
	args := m.Called(ctx, concept)
	return args.Error(0)
}

func (m *MockConceptRepo) GetByID(ctx context.Context, id uint) (*domain.PayrollConcept, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PayrollConcept), args.Error(1)
}

func (m *MockConceptRepo) GetByCode(ctx context.Context, code string) (*domain.PayrollConcept, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PayrollConcept), args.Error(1)
}

func (m *MockConceptRepo) GetActiveConcepts(ctx context.Context) ([]domain.PayrollConcept, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.PayrollConcept), args.Error(1)
}

func (m *MockConceptRepo) List(ctx context.Context, page, limit int) ([]domain.PayrollConcept, int64, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]domain.PayrollConcept), args.Get(1).(int64), args.Error(2)
}

func (m *MockConceptRepo) Update(ctx context.Context, concept *domain.PayrollConcept) error {
	args := m.Called(ctx, concept)
	return args.Error(0)
}

func (m *MockConceptRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ========================================
// Test Cases
// ========================================

func TestPayrollCalculator_Calculate_Success(t *testing.T) {
	ctx := context.Background()

	// Setup mocks
	mockPayrollRepo := new(MockPayrollRepo)
	mockPayrollItemRepo := new(MockPayrollItemRepo)
	mockEmployeeRepo := new(MockEmployeeRepo)
	mockContractRepo := new(MockContractRepo)
	mockConceptRepo := new(MockConceptRepo)

	calculator := NewPayrollCalculatorService(
		mockPayrollRepo,
		mockPayrollItemRepo,
		mockEmployeeRepo,
		mockContractRepo,
		mockConceptRepo,
	)

	// Datos de prueba
	employee := &domain.Employee{
		ID:       1,
		TenantID: 1,
		User: domain.User{
			FirstName: "Juan",
			LastName:  "Pérez",
		},
	}

	contract := &domain.EmployeeContract{
		ID:                 1,
		EmployeeID:         1,
		BaseSalary:         2000000,
		TransportAllowance: 150000,
		HousingAllowance:   300000,
	}

	concepts := []domain.PayrollConcept{
		{ID: 1, Code: domain.ConceptBaseSalary, Name: "Salario Base", Type: domain.PayrollTypeEarning, EmployeePart: 2000000},
		{ID: 2, Code: domain.ConceptTransport, Name: "Auxilio Transporte", Type: domain.PayrollTypeEarning},
		{ID: 3, Code: domain.ConceptHealth, Name: "Salud", Type: domain.PayrollTypeDeduction, Percentage: 4},
		{ID: 4, Code: domain.ConceptPension, Name: "Pensión", Type: domain.PayrollTypeDeduction, Percentage: 4},
	}

	// Configure mocks
	mockEmployeeRepo.On("GetByID", ctx, uint(1)).Return(employee, nil)
	mockContractRepo.On("GetActiveByEmployee", ctx, uint(1)).Return(contract, nil)
	mockConceptRepo.On("GetActiveConcepts", ctx).Return(concepts, nil)

	// Execute - usar exactamente 30 días para evitar ajuste de período
	req := CalculatePayrollRequest{
		EmployeeID:  1,
		PeriodStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:   time.Date(2024, 1, 30, 0, 0, 0, 0, time.UTC), // 30 días exactos
		PayDate:     time.Date(2024, 1, 30, 0, 0, 0, 0, time.UTC),
	}

	result, err := calculator.Calculate(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.Payroll.TenantID)
	assert.Equal(t, uint(1), result.Payroll.EmployeeID)
	assert.Equal(t, domain.PayrollStatusDraft, result.Payroll.Status)

	// Verificar earnings: base salary (2M del EmployeePart del concepto) + transport (150K del contract)
	// El código toma EmployeePart del concepto, no del contract para BASE_SALARY
	// Transport allowance viene del contract cuando el concepto no tiene percentage ni EmployeePart
	expectedGross := 2000000.0 + 150000.0
	assert.InDelta(t, expectedGross, result.GrossAmount, 1.0)

	// Verificar deducciones: 4% salud + 4% pension = 8% de 2000000 = 160000
	expectedDeductions := 2000000.0 * 0.08
	assert.InDelta(t, expectedDeductions, result.TotalDeductions, 1.0)

	// Net = gross - deductions
	expectedNet := expectedGross - expectedDeductions
	assert.InDelta(t, expectedNet, result.NetAmount, 1.0)

	// Verificar que se llamaron los mocks
	mockEmployeeRepo.AssertExpectations(t)
	mockContractRepo.AssertExpectations(t)
	mockConceptRepo.AssertExpectations(t)
}

func TestPayrollCalculator_Calculate_EmployeeNotFound(t *testing.T) {
	ctx := context.Background()

	mockPayrollRepo := new(MockPayrollRepo)
	mockPayrollItemRepo := new(MockPayrollItemRepo)
	mockEmployeeRepo := new(MockEmployeeRepo)
	mockContractRepo := new(MockContractRepo)
	mockConceptRepo := new(MockConceptRepo)

	calculator := NewPayrollCalculatorService(
		mockPayrollRepo,
		mockPayrollItemRepo,
		mockEmployeeRepo,
		mockContractRepo,
		mockConceptRepo,
	)

	mockEmployeeRepo.On("GetByID", ctx, uint(999)).Return(nil, domain.ErrEmployeeNotFound)

	req := CalculatePayrollRequest{
		EmployeeID:  999,
		PeriodStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
	}

	result, err := calculator.Calculate(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "employee not found")
}

func TestPayrollCalculator_Calculate_NoActiveContract(t *testing.T) {
	ctx := context.Background()

	mockPayrollRepo := new(MockPayrollRepo)
	mockPayrollItemRepo := new(MockPayrollItemRepo)
	mockEmployeeRepo := new(MockEmployeeRepo)
	mockContractRepo := new(MockContractRepo)
	mockConceptRepo := new(MockConceptRepo)

	calculator := NewPayrollCalculatorService(
		mockPayrollRepo,
		mockPayrollItemRepo,
		mockEmployeeRepo,
		mockContractRepo,
		mockConceptRepo,
	)

	employee := &domain.Employee{ID: 1, TenantID: 1}

	mockEmployeeRepo.On("GetByID", ctx, uint(1)).Return(employee, nil)
	mockContractRepo.On("GetActiveByEmployee", ctx, uint(1)).Return(nil, domain.ErrEmployeeContractNotFound)

	req := CalculatePayrollRequest{
		EmployeeID:  1,
		PeriodStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
	}

	result, err := calculator.Calculate(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no active contract")
}

func TestPayrollCalculator_Calculate_NoConcepts(t *testing.T) {
	ctx := context.Background()

	mockPayrollRepo := new(MockPayrollRepo)
	mockPayrollItemRepo := new(MockPayrollItemRepo)
	mockEmployeeRepo := new(MockEmployeeRepo)
	mockContractRepo := new(MockContractRepo)
	mockConceptRepo := new(MockConceptRepo)

	calculator := NewPayrollCalculatorService(
		mockPayrollRepo,
		mockPayrollItemRepo,
		mockEmployeeRepo,
		mockContractRepo,
		mockConceptRepo,
	)

	employee := &domain.Employee{ID: 1, TenantID: 1}
	contract := &domain.EmployeeContract{ID: 1, EmployeeID: 1, BaseSalary: 2000000}

	mockEmployeeRepo.On("GetByID", ctx, uint(1)).Return(employee, nil)
	mockContractRepo.On("GetActiveByEmployee", ctx, uint(1)).Return(contract, nil)
	mockConceptRepo.On("GetActiveConcepts", ctx).Return([]domain.PayrollConcept{}, nil)

	req := CalculatePayrollRequest{
		EmployeeID:  1,
		PeriodStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
	}

	result, err := calculator.Calculate(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no active payroll concepts")
}

func TestPayrollCalculator_Calculate_InvalidPeriod(t *testing.T) {
	ctx := context.Background()

	mockPayrollRepo := new(MockPayrollRepo)
	mockPayrollItemRepo := new(MockPayrollItemRepo)
	mockEmployeeRepo := new(MockEmployeeRepo)
	mockContractRepo := new(MockContractRepo)
	mockConceptRepo := new(MockConceptRepo)

	calculator := NewPayrollCalculatorService(
		mockPayrollRepo,
		mockPayrollItemRepo,
		mockEmployeeRepo,
		mockContractRepo,
		mockConceptRepo,
	)

	// PeriodEnd before PeriodStart
	req := CalculatePayrollRequest{
		EmployeeID:  1,
		PeriodStart: time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		PeriodEnd:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Antes que start
	}

	result, err := calculator.Calculate(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot be before")
}

func TestPayrollCalculator_Calculate_PartialPeriod(t *testing.T) {
	ctx := context.Background()

	mockPayrollRepo := new(MockPayrollRepo)
	mockPayrollItemRepo := new(MockPayrollItemRepo)
	mockEmployeeRepo := new(MockEmployeeRepo)
	mockContractRepo := new(MockContractRepo)
	mockConceptRepo := new(MockConceptRepo)

	calculator := NewPayrollCalculatorService(
		mockPayrollRepo,
		mockPayrollItemRepo,
		mockEmployeeRepo,
		mockContractRepo,
		mockConceptRepo,
	)

	employee := &domain.Employee{ID: 1, TenantID: 1}
	contract := &domain.EmployeeContract{ID: 1, EmployeeID: 1, BaseSalary: 3000000, TransportAllowance: 50000}

	// Usar percentage para que tome el baseSalary ajustado
	concepts := []domain.PayrollConcept{
		{ID: 1, Code: domain.ConceptBaseSalary, Name: "Salario Base", Type: domain.PayrollTypeEarning, Percentage: 100},
		{ID: 2, Code: domain.ConceptTransport, Name: "Auxilio Transporte", Type: domain.PayrollTypeEarning, EmployeePart: 50000},
	}

	mockEmployeeRepo.On("GetByID", ctx, uint(1)).Return(employee, nil)
	mockContractRepo.On("GetActiveByEmployee", ctx, uint(1)).Return(contract, nil)
	mockConceptRepo.On("GetActiveConcepts", ctx).Return(concepts, nil)

	// Solo 15 días de un mes de 30
	req := CalculatePayrollRequest{
		EmployeeID:  1,
		PeriodStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), // 15 días
	}

	result, err := calculator.Calculate(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Base salary ajustado: 3000000 * 15 / 30 = 1500000
	// + transport allowance fijo = 50000
	// Total = 1550000
	expectedSalary := 3000000.0*15.0/30.0 + 50000.0
	assert.InDelta(t, expectedSalary, result.GrossAmount, 1.0)
}

func TestPayrollCalculator_CalculateAndSave_NewPayroll(t *testing.T) {
	ctx := context.Background()

	mockPayrollRepo := new(MockPayrollRepo)
	mockPayrollItemRepo := new(MockPayrollItemRepo)
	mockEmployeeRepo := new(MockEmployeeRepo)
	mockContractRepo := new(MockContractRepo)
	mockConceptRepo := new(MockConceptRepo)

	calculator := NewPayrollCalculatorService(
		mockPayrollRepo,
		mockPayrollItemRepo,
		mockEmployeeRepo,
		mockContractRepo,
		mockConceptRepo,
	)

	employee := &domain.Employee{ID: 1, TenantID: 1}
	contract := &domain.EmployeeContract{ID: 1, EmployeeID: 1, BaseSalary: 2000000}
	concepts := []domain.PayrollConcept{
		{ID: 1, Code: domain.ConceptBaseSalary, Name: "Salario Base", Type: domain.PayrollTypeEarning, EmployeePart: 2000000},
	}

	mockEmployeeRepo.On("GetByID", ctx, uint(1)).Return(employee, nil)
	mockContractRepo.On("GetActiveByEmployee", ctx, uint(1)).Return(contract, nil)
	mockConceptRepo.On("GetActiveConcepts", ctx).Return(concepts, nil)
	mockPayrollRepo.On("GetByEmployeeAndPeriod", ctx, uint(1), mock.Anything, mock.Anything).Return(nil, domain.ErrPayrollNotFound)
	mockPayrollRepo.On("Create", ctx, mock.AnythingOfType("*domain.Payroll")).Return(nil)
	mockPayrollItemRepo.On("CreateBatch", ctx, mock.AnythingOfType("[]domain.PayrollItem")).Return(nil)

	req := CalculatePayrollRequest{
		EmployeeID:  1,
		PeriodStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		PayDate:     time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
	}

	result, err := calculator.CalculateAndSave(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotZero(t, result.Payroll.ID)

	mockPayrollRepo.AssertCalled(t, "Create", ctx, mock.AnythingOfType("*domain.Payroll"))
	mockPayrollItemRepo.AssertCalled(t, "CreateBatch", ctx, mock.AnythingOfType("[]domain.PayrollItem"))
}

func TestPayrollCalculator_CalculateAndSave_UpdateExisting(t *testing.T) {
	ctx := context.Background()

	mockPayrollRepo := new(MockPayrollRepo)
	mockPayrollItemRepo := new(MockPayrollItemRepo)
	mockEmployeeRepo := new(MockEmployeeRepo)
	mockContractRepo := new(MockContractRepo)
	mockConceptRepo := new(MockConceptRepo)

	calculator := NewPayrollCalculatorService(
		mockPayrollRepo,
		mockPayrollItemRepo,
		mockEmployeeRepo,
		mockContractRepo,
		mockConceptRepo,
	)

	employee := &domain.Employee{ID: 1, TenantID: 1}
	contract := &domain.EmployeeContract{ID: 1, EmployeeID: 1, BaseSalary: 2000000}
	concepts := []domain.PayrollConcept{
		{ID: 1, Code: domain.ConceptBaseSalary, Name: "Salario Base", Type: domain.PayrollTypeEarning, EmployeePart: 2000000},
	}

	existingPayroll := &domain.Payroll{ID: 5, EmployeeID: 1, Status: domain.PayrollStatusDraft}

	mockEmployeeRepo.On("GetByID", ctx, uint(1)).Return(employee, nil)
	mockContractRepo.On("GetActiveByEmployee", ctx, uint(1)).Return(contract, nil)
	mockConceptRepo.On("GetActiveConcepts", ctx).Return(concepts, nil)
	mockPayrollRepo.On("GetByEmployeeAndPeriod", ctx, uint(1), mock.Anything, mock.Anything).Return(existingPayroll, nil)
	mockPayrollRepo.On("Update", ctx, mock.AnythingOfType("*domain.Payroll")).Return(nil)
	mockPayrollItemRepo.On("DeleteByPayrollID", ctx, uint(5)).Return(nil)
	mockPayrollItemRepo.On("CreateBatch", ctx, mock.AnythingOfType("[]domain.PayrollItem")).Return(nil)

	req := CalculatePayrollRequest{
		EmployeeID:  1,
		PeriodStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
	}

	result, err := calculator.CalculateAndSave(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(5), result.Payroll.ID) // ID del existente

	// Verificar que se actualizó en vez de crear
	mockPayrollRepo.AssertCalled(t, "Update", ctx, mock.AnythingOfType("*domain.Payroll"))
	mockPayrollItemRepo.AssertCalled(t, "DeleteByPayrollID", ctx, uint(5))
}

// ========================================
// Test para PayrollStateService
// ========================================

type MockPaymentRepo struct {
	mock.Mock
}

func (m *MockPaymentRepo) Create(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	if payment != nil {
		payment.ID = 1
	}
	return args.Error(0)
}

func (m *MockPaymentRepo) GetByID(ctx context.Context, id uint) (*domain.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepo) GetByPayrollID(ctx context.Context, payrollID uint) (*domain.Payment, error) {
	args := m.Called(ctx, payrollID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestPayrollStateService_MarkAsPaid_Success(t *testing.T) {
	ctx := context.Background()

	mockPayrollRepo := new(MockPayrollRepo)
	mockPaymentRepo := new(MockPaymentRepo)
	mockEmployeeRepo := new(MockEmployeeRepo)

	stateSvc := NewPayrollStateService(mockPayrollRepo, mockPaymentRepo, mockEmployeeRepo)

	payroll := &domain.Payroll{
		ID:         1,
		EmployeeID: 1,
		Status:     domain.PayrollStatusCalculated,
		NetAmount:  1800000,
	}

	mockPayrollRepo.On("GetByID", ctx, uint(1)).Return(payroll, nil)
	mockPaymentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil)
	mockPayrollRepo.On("Update", ctx, mock.AnythingOfType("*domain.Payroll")).Return(nil)

	payment, err := stateSvc.MarkAsPaid(ctx, 1, "bank_transfer")

	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, float64(1800000), payment.Amount)
	assert.Equal(t, "completed", payment.Status)

	// Verificar que el estado cambió
	mockPayrollRepo.AssertCalled(t, "Update", ctx, mock.MatchedBy(func(p *domain.Payroll) bool {
		return p.Status == domain.PayrollStatusPaid
	}))
}

func TestPayrollStateService_MarkAsPaid_AlreadyPaid(t *testing.T) {
	ctx := context.Background()

	mockPayrollRepo := new(MockPayrollRepo)
	mockPaymentRepo := new(MockPaymentRepo)
	mockEmployeeRepo := new(MockEmployeeRepo)

	stateSvc := NewPayrollStateService(mockPayrollRepo, mockPaymentRepo, mockEmployeeRepo)

	payroll := &domain.Payroll{
		ID:     1,
		Status: domain.PayrollStatusPaid,
	}

	mockPayrollRepo.On("GetByID", ctx, uint(1)).Return(payroll, nil)

	_, err := stateSvc.MarkAsPaid(ctx, 1, "bank_transfer")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already paid")
}

func TestPayrollStateService_RevertToDraft_Success(t *testing.T) {
	ctx := context.Background()

	mockPayrollRepo := new(MockPayrollRepo)
	mockPaymentRepo := new(MockPaymentRepo)
	mockEmployeeRepo := new(MockEmployeeRepo)

	stateSvc := NewPayrollStateService(mockPayrollRepo, mockPaymentRepo, mockEmployeeRepo)

	payroll := &domain.Payroll{
		ID:     1,
		Status: domain.PayrollStatusCalculated,
	}

	mockPayrollRepo.On("GetByID", ctx, uint(1)).Return(payroll, nil)
	mockPayrollRepo.On("Update", ctx, mock.AnythingOfType("*domain.Payroll")).Return(nil)

	err := stateSvc.RevertToDraft(ctx, 1)

	assert.NoError(t, err)
	mockPayrollRepo.AssertCalled(t, "Update", ctx, mock.MatchedBy(func(p *domain.Payroll) bool {
		return p.Status == domain.PayrollStatusDraft
	}))
}

func TestPayrollStateService_RevertToDraft_CannotRevertPaid(t *testing.T) {
	ctx := context.Background()

	mockPayrollRepo := new(MockPayrollRepo)
	mockPaymentRepo := new(MockPaymentRepo)
	mockEmployeeRepo := new(MockEmployeeRepo)

	stateSvc := NewPayrollStateService(mockPayrollRepo, mockPaymentRepo, mockEmployeeRepo)

	payroll := &domain.Payroll{
		ID:     1,
		Status: domain.PayrollStatusPaid,
	}

	mockPayrollRepo.On("GetByID", ctx, uint(1)).Return(payroll, nil)

	err := stateSvc.RevertToDraft(ctx, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot revert a paid payroll")
}
