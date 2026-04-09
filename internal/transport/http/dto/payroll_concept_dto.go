package dto

import (
	"strings"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
)

// ========================================
// PayrollConcept DTOs
// ========================================

// CreatePayrollConceptRequest representa el DTO para crear conceptos de nómina
type CreatePayrollConceptRequest struct {
	Code         string  `json:"code" binding:"required,max=30"`
	Name         string  `json:"name" binding:"required,max=100"`
	Type         string  `json:"type" binding:"required,oneof=earning deduction employer_contribution"`
	Description  string  `json:"description,omitempty" binding:"max=255"`
	Percentage   float64 `json:"percentage,omitempty" binding:"min=0,max=100"`
	EmployeePart float64 `json:"employee_part,omitempty" binding:"min=0,max=100"`
	EmployerPart float64 `json:"employer_part,omitempty" binding:"min=0,max=100"`
	IsMandatory  *bool   `json:"is_mandatory,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

// UpdatePayrollConceptRequest representa el DTO para actualizar conceptos
type UpdatePayrollConceptRequest struct {
	Name         *string  `json:"name,omitempty"`
	Type         *string  `json:"type,omitempty"`
	Description  *string  `json:"description,omitempty"`
	Percentage   *float64 `json:"percentage,omitempty"`
	EmployeePart *float64 `json:"employee_part,omitempty"`
	EmployerPart *float64 `json:"employer_part,omitempty"`
	IsMandatory  *bool    `json:"is_mandatory,omitempty"`
	IsActive     *bool    `json:"is_active,omitempty"`
}

// PayrollConceptResponse representa la respuesta de un concepto de nómina
type PayrollConceptResponse struct {
	ID           uint      `json:"id"`
	TenantID     uint      `json:"tenant_id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Description  string    `json:"description,omitempty"`
	Percentage   float64   `json:"percentage"`
	EmployeePart float64   `json:"employee_part"`
	EmployerPart float64   `json:"employer_part"`
	IsMandatory  bool      `json:"is_mandatory"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ========================================
// DTO Conversion Methods
// ========================================

// ToDomain convierte CreatePayrollConceptRequest a domain.PayrollConcept
func (r *CreatePayrollConceptRequest) ToDomain() *domain.PayrollConcept {
	concept := &domain.PayrollConcept{
		Code:         strings.ToUpper(strings.TrimSpace(r.Code)),
		Name:         strings.TrimSpace(r.Name),
		Type:         r.Type,
		Description:  r.Description,
		Percentage:   r.Percentage,
		EmployeePart: r.EmployeePart,
		EmployerPart: r.EmployerPart,
	}

	if r.IsMandatory != nil {
		concept.IsMandatory = *r.IsMandatory
	} else {
		concept.IsMandatory = false
	}

	if r.IsActive != nil {
		concept.IsActive = *r.IsActive
	} else {
		concept.IsActive = true
	}

	return concept
}

// ToResponse convierte domain.PayrollConcept a PayrollConceptResponse
func ToPayrollConceptResponse(concept *domain.PayrollConcept) *PayrollConceptResponse {
	resp := &PayrollConceptResponse{
		ID:           concept.ID,
		TenantID:     concept.TenantID,
		Code:         concept.Code,
		Name:         concept.Name,
		Type:         concept.Type,
		Description:  concept.Description,
		Percentage:   concept.Percentage,
		EmployeePart: concept.EmployeePart,
		EmployerPart: concept.EmployerPart,
		IsMandatory:  concept.IsMandatory,
		IsActive:     concept.IsActive,
		CreatedAt:    concept.CreatedAt,
		UpdatedAt:    concept.UpdatedAt,
	}
	return resp
}

// Sanitize limpia los datos del DTO
func (r *CreatePayrollConceptRequest) Sanitize() {
	r.Code = strings.ToUpper(strings.TrimSpace(r.Code))
	r.Name = strings.TrimSpace(r.Name)
	if r.Description != "" {
		r.Description = strings.TrimSpace(r.Description)
	}
}
