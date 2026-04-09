package dto

import (
	"strings"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
)

// ========================================
// Employee DTOs
// ========================================

// CreateEmployeeRequest representa el DTO para crear empleados
type CreateEmployeeRequest struct {
	UserID       uint  `json:"user_id" binding:"required,min=1"`
	DepartmentID *uint `json:"department_id,omitempty"`
	PositionID   *uint `json:"position_id,omitempty"`
	IsActive     *bool `json:"is_active,omitempty"`
}

// UpdateEmployeeRequest representa el DTO para actualizar empleados
type UpdateEmployeeRequest struct {
	DepartmentID *uint `json:"department_id,omitempty"`
	PositionID   *uint `json:"position_id,omitempty"`
	IsActive     *bool `json:"is_active,omitempty"`
}

// EmployeeResponse representa la respuesta de un empleado
type EmployeeResponse struct {
	ID           uint                `json:"id"`
	TenantID     uint                `json:"tenant_id"`
	UserID       uint                `json:"user_id"`
	DepartmentID *uint               `json:"department_id,omitempty"`
	PositionID   *uint               `json:"position_id,omitempty"`
	IsActive     bool                `json:"is_active"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	User         *UserResponse       `json:"user,omitempty"`
	Department   *DepartmentResponse `json:"department,omitempty"`
	Position     *PositionResponse   `json:"position,omitempty"`
	Contracts    []ContractResponse  `json:"contracts,omitempty"`
}

// UserResponse representa la respuesta de usuario embebido
type UserResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Dni       string `json:"dni"`
	Gender    string `json:"gender"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	BirthDay  string `json:"birth_day"`
}

// DepartmentResponse representa la respuesta de departamento
type DepartmentResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// PositionResponse representa la respuesta de posición
type PositionResponse struct {
	ID           uint   `json:"id"`
	NamePosition string `json:"name_position"`
	Description  string `json:"description"`
}

// ContractResponse representa la respuesta de contrato
type ContractResponse struct {
	ID                  uint    `json:"id"`
	ContractTypeID      uint    `json:"contract_type_id,omitempty"`
	BaseSalary          float64 `json:"base_salary"`
	Currency            string  `json:"currency"`
	StartDate           string  `json:"start_date"`
	EndDate             *string `json:"end_date,omitempty"`
	IsActive            bool    `json:"is_active"`
	WorkHoursPerDay     float64 `json:"work_hours_per_day"`
	WorkDaysPerWeek     float64 `json:"work_days_per_week"`
	HealthContribution  float64 `json:"health_contribution"`
	PensionContribution float64 `json:"pension_contribution"`
	TransportAllowance  float64 `json:"transport_allowance"`
	HousingAllowance    float64 `json:"housing_allowance"`
}

// ========================================
// Employee Contract DTOs
// ========================================

// CreateEmployeeContractRequest representa el DTO para crear contratos
type CreateEmployeeContractRequest struct {
	EmployeeID          uint    `json:"employee_id" binding:"required,min=1"`
	ContractTypeID      *uint   `json:"contract_type_id,omitempty"`
	BaseSalary          float64 `json:"base_salary" binding:"required,min=0"`
	Currency            string  `json:"currency" binding:"required,len=3"`
	StartDate           string  `json:"start_date" binding:"required"`
	EndDate             *string `json:"end_date,omitempty"`
	WorkHoursPerDay     float64 `json:"work_hours_per_day" binding:"min=0,max=24"`
	WorkDaysPerWeek     float64 `json:"work_days_per_week" binding:"min=0,max=7"`
	HealthContribution  float64 `json:"health_contribution" binding:"min=0,max=100"`
	PensionContribution float64 `json:"pension_contribution" binding:"min=0,max=100"`
	TransportAllowance  float64 `json:"transport_allowance" binding:"min=0"`
	HousingAllowance    float64 `json:"housing_allowance" binding:"min=0"`
}

// ========================================
// DTO Conversion Methods
// ========================================

// ToDomain convierte CreateEmployeeRequest a domain.Employee
func (r *CreateEmployeeRequest) ToDomain() *domain.Employee {
	emp := &domain.Employee{
		UserID: r.UserID,
	}

	if r.DepartmentID != nil {
		emp.DepartmentID = *r.DepartmentID
	}
	if r.PositionID != nil {
		emp.PositionID = *r.PositionID
	}
	if r.IsActive != nil {
		emp.IsActive = *r.IsActive
	} else {
		emp.IsActive = true
	}

	return emp
}

// ToResponse convierte domain.Employee a EmployeeResponse
func ToEmployeeResponse(emp *domain.Employee) *EmployeeResponse {
	resp := &EmployeeResponse{
		ID:           emp.ID,
		TenantID:     emp.TenantID,
		UserID:       emp.UserID,
		DepartmentID: nil,
		PositionID:   nil,
		IsActive:     emp.IsActive,
		CreatedAt:    emp.CreatedAt,
		UpdatedAt:    emp.UpdatedAt,
	}

	if emp.DepartmentID > 0 {
		resp.DepartmentID = &emp.DepartmentID
	}
	if emp.PositionID > 0 {
		resp.PositionID = &emp.PositionID
	}

	// Include user if preloaded
	if emp.User.ID > 0 {
		resp.User = &UserResponse{
			ID:        emp.User.ID,
			FirstName: emp.User.FirstName,
			LastName:  emp.User.LastName,
			Dni:       emp.User.Dni,
			Gender:    emp.User.Gender,
			Phone:     emp.User.Phone,
			Email:     emp.User.Email,
			BirthDay:  emp.User.BirthDay.Format("2006-01-02"),
		}
	}

	// Include department if preloaded
	if emp.Department.ID > 0 {
		resp.Department = &DepartmentResponse{
			ID:   emp.Department.ID,
			Name: emp.Department.Name,
			Code: emp.Department.Code,
		}
	}

	// Include position if preloaded
	if emp.Position.ID > 0 {
		resp.Position = &PositionResponse{
			ID:           emp.Position.ID,
			NamePosition: emp.Position.NamePosition,
			Description:  emp.Position.Description,
		}
	}

	// Include contracts if preloaded
	if len(emp.Contracts) > 0 {
		resp.Contracts = make([]ContractResponse, len(emp.Contracts))
		for i, c := range emp.Contracts {
			contractResp := ContractResponse{
				ID:                  c.ID,
				BaseSalary:          c.BaseSalary,
				Currency:            c.Currency,
				StartDate:           c.StartDate.Format("2006-01-02"),
				IsActive:            c.IsActive,
				WorkHoursPerDay:     c.WorkHoursPerDay,
				WorkDaysPerWeek:     c.WorkDaysPerWeek,
				HealthContribution:  c.HealthContribution,
				PensionContribution: c.PensionContribution,
				TransportAllowance:  c.TransportAllowance,
				HousingAllowance:    c.HousingAllowance,
			}
			if c.ContractTypeID > 0 {
				contractResp.ContractTypeID = c.ContractTypeID
			}
			if c.EndDate != nil {
				endDate := c.EndDate.Format("2006-01-02")
				contractResp.EndDate = &endDate
			}
			resp.Contracts[i] = contractResp
		}
	}

	return resp
}

// Sanitize limpia los datos del DTO
func (r *CreateEmployeeRequest) Sanitize() {
	// Currently no fields to sanitize
}

// ToDomain convierte CreateEmployeeContractRequest a domain.EmployeeContract
func (r *CreateEmployeeContractRequest) ToDomain() (*domain.EmployeeContract, error) {
	startDate, err := time.Parse("2006-01-02", r.StartDate)
	if err != nil {
		return nil, err
	}

	contract := &domain.EmployeeContract{
		EmployeeID:          r.EmployeeID,
		BaseSalary:          r.BaseSalary,
		Currency:            strings.ToUpper(r.Currency),
		StartDate:           startDate,
		WorkHoursPerDay:     r.WorkHoursPerDay,
		WorkDaysPerWeek:     r.WorkDaysPerWeek,
		HealthContribution:  r.HealthContribution,
		PensionContribution: r.PensionContribution,
		TransportAllowance:  r.TransportAllowance,
		HousingAllowance:    r.HousingAllowance,
		IsActive:            true,
	}

	if r.ContractTypeID != nil {
		contract.ContractTypeID = *r.ContractTypeID
	}

	if r.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *r.EndDate)
		if err != nil {
			return nil, err
		}
		contract.EndDate = &endDate
	}

	return contract, nil
}
