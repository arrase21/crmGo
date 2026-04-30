package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/arrase21/crm-users/internal/domain"
)

// ========================================
// Errores específicos del servicio
// ========================================

var (
	ErrEmployeeNotFound         = errors.New("employee not found")
	ErrUserNotFound             = errors.New("user not found")
	ErrDuplicateEmployee        = errors.New("user already has an active employee record")
	ErrCannotActivateEmployee   = errors.New("cannot activate employee")
	ErrCannotDeactivateEmployee = errors.New("cannot deactivate employee: has active contract")
)

// ========================================
// DTOs para requests/responses
// ========================================

type CreateEmployeeRequest struct {
	UserID       uint   `json:"user_id" validate:"required"`
	DepartmentID *uint  `json:"department_id"`
	PositionID   *uint  `json:"position_id"`
	EmployeeCode string `json:"employee_code"`
	Notes        string `json:"notes"`
}

type UpdateEmployeeRequest struct {
	DepartmentID *uint   `json:"department_id"`
	PositionID   *uint   `json:"position_id"`
	EmployeeCode *string `json:"employee_code"`
	Notes        *string `json:"notes"`
}

type ChangeEmployeeStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active inactive suspended"`
	Reason string `json:"reason"`
}

// EmployeeResponse es la respuesta pública del empleado
type EmployeeResponse struct {
	ID             uint             `json:"id"`
	EmployeeCode   string           `json:"employee_code"`
	Status         string           `json:"status"`
	IsActive       bool             `json:"is_active"`
	Notes          string           `json:"notes"`
	CreatedAt      string           `json:"created_at"`
	UpdatedAt      string           `json:"updated_at"`
	User           UserSummary      `json:"user"`
	Department     *SimpleEntity    `json:"department,omitempty"`
	Position       *SimpleEntity    `json:"position,omitempty"`
	ActiveContract *ContractSummary `json:"active_contract,omitempty"`
}

type UserSummary struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	DNI       string `json:"dni"`
	Email     string `json:"email"`
}

type SimpleEntity struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ContractSummary struct {
	ID           uint    `json:"id"`
	ContractType string  `json:"contract_type"`
	BaseSalary   float64 `json:"base_salary"`
	Currency     string  `json:"currency"`
	StartDate    string  `json:"start_date"`
	EndDate      *string `json:"end_date,omitempty"`
}

// ========================================
// Service principal
// ========================================

type EmployeeService struct {
	employeeRepo domain.EmployeeRepo
	userRepo     domain.UserRepo
}

func NewEmployeeService(empRepo domain.EmployeeRepo, userRepo domain.UserRepo) *EmployeeService {
	return &EmployeeService{
		employeeRepo: empRepo,
		userRepo:     userRepo,
	}
}

// ========================================
// Métodos de negocio
// ========================================

// Create crea un nuevo empleado con validaciones de negocio
func (s *EmployeeService) Create(ctx context.Context, req CreateEmployeeRequest) (*EmployeeResponse, error) {
	// 1. Validar que el usuario existe
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	// 2. Verificar que el usuario no tenga già un empleado
	exists, err := s.employeeRepo.ExistsByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("error checking employee existence: %w", err)
	}
	if exists {
		return nil, ErrDuplicateEmployee
	}

	// 3. Obtener el TenantID del contexto
	tenantID, err := s.extractTenantID(ctx)
	if err != nil {
		return nil, err
	}

	// 4. Crear el empleado
	employee := &domain.Employee{
		TenantID:     tenantID,
		UserID:       req.UserID,
		EmployeeCode: req.EmployeeCode,
		Notes:        req.Notes,
		Status:       domain.EmployeeStatusActive,
		IsActive:     true,
	}

	if req.DepartmentID != nil {
		employee.DepartmentID = *req.DepartmentID
	}
	if req.PositionID != nil {
		employee.PositionID = *req.PositionID
	}

	// El modelo valida su propio estado
	if err := employee.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if err := s.employeeRepo.Create(ctx, employee); err != nil {
		return nil, fmt.Errorf("error creating employee: %w", err)
	}

	// 5. Cargar relaciones para la respuesta
	employee.User = *user
	return s.toResponse(employee), nil
}

// GetByID obtiene un empleado por ID
func (s *EmployeeService) GetByID(ctx context.Context, id uint) (*EmployeeResponse, error) {
	if id == 0 {
		return nil, errors.New("invalid employee id")
	}

	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("error fetching employee: %w", err)
	}

	return s.toResponse(employee), nil
}

// GetByUserID obtiene un empleado por UserID
func (s *EmployeeService) GetByUserID(ctx context.Context, userID uint) (*EmployeeResponse, error) {
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}

	employee, err := s.employeeRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("error fetching employee: %w", err)
	}

	return s.toResponse(employee), nil
}

// List lista empleados (mantiene compatibilidad con page/limit)
func (s *EmployeeService) List(ctx context.Context, page, limit int) ([]EmployeeResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	employees, total, err := s.employeeRepo.List(ctx, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing employees: %w", err)
	}

	response := make([]EmployeeResponse, len(employees))
	for i, emp := range employees {
		response[i] = *s.toResponse(&emp)
	}

	return response, total, nil
}

// Update actualiza un empleado
func (s *EmployeeService) Update(ctx context.Context, id uint, req UpdateEmployeeRequest) (*EmployeeResponse, error) {
	if id == 0 {
		return nil, errors.New("invalid employee id")
	}

	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("error fetching employee: %w", err)
	}

	// Aplicar solo los campos no nil
	if req.DepartmentID != nil {
		employee.DepartmentID = *req.DepartmentID
	}
	if req.PositionID != nil {
		employee.PositionID = *req.PositionID
	}
	if req.EmployeeCode != nil {
		employee.EmployeeCode = *req.EmployeeCode
	}
	if req.Notes != nil {
		employee.Notes = *req.Notes
	}

	if err := employee.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if err := s.employeeRepo.Update(ctx, employee); err != nil {
		return nil, fmt.Errorf("error updating employee: %w", err)
	}

	// Recargar para obtener relaciones actualizadas
	employee, err = s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error reloading employee: %w", err)
	}

	return s.toResponse(employee), nil
}

// ChangeStatus cambia el estado del empleado
func (s *EmployeeService) ChangeStatus(ctx context.Context, id uint, req ChangeEmployeeStatusRequest) (*EmployeeResponse, error) {
	if id == 0 {
		return nil, errors.New("invalid employee id")
	}

	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("error fetching employee: %w", err)
	}

	// Aplicar el cambio de estado (el modelo valida las transiciones)
	if err := employee.SetStatus(req.Status); err != nil {
		return nil, err
	}

	if err := s.employeeRepo.Update(ctx, employee); err != nil {
		return nil, fmt.Errorf("error updating employee status: %w", err)
	}

	employee, err = s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error reloading employee: %w", err)
	}

	return s.toResponse(employee), nil
}

// Deactivate desactiva un empleado (soft delete con validaciones)
func (s *EmployeeService) Deactivate(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid employee id")
	}

	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			return ErrEmployeeNotFound
		}
		return fmt.Errorf("error fetching employee: %w", err)
	}

	// Validar que puede ser desactivado
	if !employee.CanBeDeactivated() {
		return ErrCannotDeactivateEmployee
	}

	employee.IsActive = false
	employee.Status = domain.EmployeeStatusInactive

	return s.employeeRepo.Update(ctx, employee)
}

// GetStatistics retorna estadísticas de empleados
func (s *EmployeeService) GetStatistics(ctx context.Context) (domain.EmployeeStatistics, error) {
	tenantID, err := s.extractTenantID(ctx)
	if err != nil {
		return domain.EmployeeStatistics{}, err
	}

	return s.employeeRepo.GetStatistics(ctx, tenantID)
}

// ========================================
// Métodos helper privados
// ========================================

func (s *EmployeeService) toResponse(emp *domain.Employee) *EmployeeResponse {
	resp := &EmployeeResponse{
		ID:           emp.ID,
		EmployeeCode: emp.EmployeeCode,
		Status:       emp.Status,
		IsActive:     emp.IsActive,
		Notes:        emp.Notes,
		CreatedAt:    emp.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    emp.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		User: UserSummary{
			ID:        emp.User.ID,
			FirstName: emp.User.FirstName,
			LastName:  emp.User.LastName,
			DNI:       emp.User.Dni,
			Email:     emp.User.Email,
		},
	}

	if emp.DepartmentID != 0 {
		resp.Department = &SimpleEntity{
			ID:   emp.Department.ID,
			Name: emp.Department.Name,
		}
	}

	if emp.PositionID != 0 {
		resp.Position = &SimpleEntity{
			ID:   emp.Position.ID,
			Name: emp.Position.NamePosition,
		}
	}

	contract := emp.GetMainContract()
	if contract != nil {
		resp.ActiveContract = &ContractSummary{
			ID:           contract.ID,
			ContractType: contract.ContractType.Name,
			BaseSalary:   contract.BaseSalary,
			Currency:     contract.Currency,
			StartDate:    contract.StartDate.Format("2006-01-02"),
		}
		if contract.EndDate != nil {
			endDate := contract.EndDate.Format("2006-01-02")
			resp.ActiveContract.EndDate = &endDate
		}
	}

	return resp
}

func (s *EmployeeService) extractTenantID(ctx context.Context) (uint, error) {
	tenantID, ok := ctx.Value(domain.TenantIDKey).(uint)
	if !ok || tenantID == 0 {
		return 0, errors.New("tenant not found in context")
	}
	return tenantID, nil
}
