package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
)

// PayrollBatchService procesa nóminas de múltiples empleados
type PayrollBatchService struct {
	payrollRepo     domain.PayrollRepo
	payrollItemRepo domain.PayrollItemRepo
	employeeRepo    domain.EmployeeRepo
	contractRepo    domain.EmployeeContractRepo
	conceptRepo     domain.PayrollConceptRepo
	stateService    *PayrollStateService
}

func NewPayrollBatchService(
	payrollRepo domain.PayrollRepo,
	payrollItemRepo domain.PayrollItemRepo,
	employeeRepo domain.EmployeeRepo,
	contractRepo domain.EmployeeContractRepo,
	conceptRepo domain.PayrollConceptRepo,
	stateService *PayrollStateService,
) *PayrollBatchService {
	return &PayrollBatchService{
		payrollRepo:     payrollRepo,
		payrollItemRepo: payrollItemRepo,
		employeeRepo:    employeeRepo,
		contractRepo:    contractRepo,
		conceptRepo:     conceptRepo,
		stateService:    stateService,
	}
}

// BatchResult representa el resultado del procesamiento de un empleado
type BatchResult struct {
	EmployeeID   uint    `json:"employee_id"`
	EmployeeName string  `json:"employee_name"`
	PayrollID    uint    `json:"payroll_id,omitempty"`
	NetAmount    float64 `json:"net_amount"`
	Success      bool    `json:"success"`
	Error        string  `json:"error,omitempty"`
}

// BatchPayrollRequest representa la solicitud de procesamiento batch
type BatchPayrollRequest struct {
	PeriodStart time.Time `json:"period_start" binding:"required"`
	PeriodEnd   time.Time `json:"period_end" binding:"required"`
	PayDate     time.Time `json:"pay_date"`
	// Procesamiento paralelo
	Concurrent bool `json:"concurrent"`
	// Limitar a ciertos departamentos (vacío = todos)
	DepartmentIDs []uint `json:"department_ids,omitempty"`
}

// BatchPayrollResponse representa la respuesta del procesamiento batch
type BatchPayrollResponse struct {
	PeriodStart    string        `json:"period_start"`
	PeriodEnd      string        `json:"period_end"`
	TotalEmployees int           `json:"total_employees"`
	Successful     int           `json:"successful"`
	Failed         int           `json:"failed"`
	TotalNet       float64       `json:"total_net_amount"`
	Results        []BatchResult `json:"results"`
	ProcessedAt    string        `json:"processed_at"`
}

// CalculatePeriodSummary procesa nóminas de todos los empleados activos
func (s *PayrollBatchService) CalculatePeriodSummary(ctx context.Context, req BatchPayrollRequest) (*BatchPayrollResponse, error) {
	// Validaciones
	if req.PeriodStart.IsZero() || req.PeriodEnd.IsZero() {
		return nil, domain.ErrInvalidPeriod
	}
	if req.PeriodEnd.Before(req.PeriodStart) {
		return nil, domain.ErrInvalidPeriod
	}

	// Si no hay pay_date, usar period end
	if req.PayDate.IsZero() {
		req.PayDate = req.PeriodEnd
	}

	// Obtener empleados activos
	employees, total, err := s.employeeRepo.ListActive(ctx, 1, 10000) // Paginar si hay muchos
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return nil, errors.New("no active employees found")
	}

	// Filtrar por departamentos si se especificó
	if len(req.DepartmentIDs) > 0 {
		employees = s.filterByDepartments(employees, req.DepartmentIDs)
	}

	response := &BatchPayrollResponse{
		PeriodStart:    req.PeriodStart.Format("2006-01-02"),
		PeriodEnd:      req.PeriodEnd.Format("2006-01-02"),
		TotalEmployees: len(employees),
		Results:        make([]BatchResult, 0, len(employees)),
		ProcessedAt:    time.Now().Format(time.RFC3339),
	}

	// Procesar según el modo
	if req.Concurrent {
		s.processConcurrently(ctx, employees, req, response)
	} else {
		s.processSequentially(ctx, employees, req, response)
	}

	return response, nil
}

// processSequentially procesa empleado por empleado
func (s *PayrollBatchService) processSequentially(
	ctx context.Context,
	employees []domain.Employee,
	req BatchPayrollRequest,
	response *BatchPayrollResponse,
) {
	for _, emp := range employees {
		result := s.processEmployee(ctx, emp, req)
		response.Results = append(response.Results, result)

		if result.Success {
			response.Successful++
			response.TotalNet += result.NetAmount
		} else {
			response.Failed++
		}
	}
}

// processConcurrently procesa en paralelo usando goroutines
func (s *PayrollBatchService) processConcurrently(
	ctx context.Context,
	employees []domain.Employee,
	req BatchPayrollRequest,
	response *BatchPayrollResponse,
) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Canal para colectar resultados
	results := make(chan BatchResult, len(employees))

	for _, emp := range employees {
		wg.Add(1)
		go func(employee domain.Employee) {
			defer wg.Done()
			result := s.processEmployee(ctx, employee, req)
			results <- result
		}(emp)
	}

	// Goroutine para cerrar el canal cuando termine
	go func() {
		wg.Wait()
		close(results)
	}()

	// Colectar resultados
	for result := range results {
		mu.Lock()
		response.Results = append(response.Results, result)
		if result.Success {
			response.Successful++
			response.TotalNet += result.NetAmount
		} else {
			response.Failed++
		}
		mu.Unlock()
	}
}

// processEmployee procesa la nómina de un empleado individual
func (s *PayrollBatchService) processEmployee(
	ctx context.Context,
	employee domain.Employee,
	req BatchPayrollRequest,
) BatchResult {
	result := BatchResult{
		EmployeeID: employee.ID,
	}

	// Obtener nombre del empleado
	if employee.User.FirstName != "" {
		result.EmployeeName = employee.User.FirstName + " " + employee.User.LastName
	}

	// Validar que no exista ya una nómina pagada para este periodo
	if err := s.stateService.ValidatePayrollCreation(ctx, employee.ID, req.PeriodStart, req.PeriodEnd); err != nil {
		result.Error = err.Error()
		return result
	}

	// Calcular y guardar
	calcReq := CalculatePayrollRequest{
		EmployeeID:  employee.ID,
		PeriodStart: req.PeriodStart,
		PeriodEnd:   req.PeriodEnd,
		PayDate:     req.PayDate,
	}

	calculator := NewPayrollCalculatorService(
		s.payrollRepo,
		s.payrollItemRepo,
		s.employeeRepo,
		s.contractRepo,
		s.conceptRepo,
	)

	calculated, err := calculator.CalculateAndSave(ctx, calcReq)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	// Marcar como calculada
	if err := s.stateService.MarkAsCalculated(ctx, calculated.Payroll.ID); err != nil {
		result.Error = "payroll saved but could not mark as calculated: " + err.Error()
		return result
	}

	result.PayrollID = calculated.Payroll.ID
	result.NetAmount = calculated.NetAmount
	result.Success = true

	return result
}

// filterByDepartments filtra empleados por departamentos
func (s *PayrollBatchService) filterByDepartments(employees []domain.Employee, deptIDs []uint) []domain.Employee {
	deptMap := make(map[uint]bool)
	for _, id := range deptIDs {
		deptMap[id] = true
	}

	filtered := make([]domain.Employee, 0)
	for _, emp := range employees {
		if deptMap[emp.DepartmentID] {
			filtered = append(filtered, emp)
		}
	}
	return filtered
}
