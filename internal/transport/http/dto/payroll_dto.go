package dto

import (
	"time"

	"github.com/arrase21/crm-users/internal/domain"
	"github.com/arrase21/crm-users/internal/service"
)

// CalculatePayrollRequest representa el request para calcular una nómina
type CalculatePayrollRequest struct {
	EmployeeID  uint      `json:"employee_id" binding:"required"`
	PeriodStart time.Time `json:"period_start" binding:"required"`
	PeriodEnd   time.Time `json:"period_end" binding:"required"`
	PayDate     time.Time `json:"pay_date"`
}

// PayrollItemResponse representa un item de nómina en la respuesta
type PayrollItemResponse struct {
	ID           uint    `json:"id"`
	ConceptID    uint    `json:"concept_id"`
	Type         string  `json:"type"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Amount       float64 `json:"amount"`
	CalculatedAt string  `json:"calculated_at"`
}

// PayrollResponse representa la respuesta de una nómina
type PayrollResponse struct {
	ID              uint                  `json:"id"`
	TenantID        uint                  `json:"tenant_id"`
	EmployeeID      uint                  `json:"employee_id"`
	EmployeeName    string                `json:"employee_name"`
	PeriodStart     string                `json:"period_start"`
	PeriodEnd       string                `json:"period_end"`
	PayDate         string                `json:"pay_date"`
	GrossAmount     float64               `json:"gross_amount"`
	TotalDeductions float64               `json:"total_deductions"`
	NetAmount       float64               `json:"net_amount"`
	Status          string                `json:"status"`
	Items           []PayrollItemResponse `json:"items"`
	CreatedAt       string                `json:"created_at"`
	UpdatedAt       string                `json:"updated_at"`
}

// CalculatedPayrollResponse representa la respuesta del cálculo de nómina
type CalculatedPayrollResponse struct {
	Payroll         PayrollResponse `json:"payroll"`
	GrossAmount     float64         `json:"gross_amount"`
	TotalDeductions float64         `json:"total_deductions"`
	NetAmount       float64         `json:"net_amount"`
	Message         string          `json:"message"`
}

// ToCalculatedPayrollResponse convierte el resultado del cálculo a DTO
func ToCalculatedPayrollResponse(calc *service.CalculatedPayroll) *CalculatedPayrollResponse {
	items := make([]PayrollItemResponse, len(calc.Items))
	for i, item := range calc.Items {
		items[i] = PayrollItemResponse{
			ID:           item.ID,
			ConceptID:    item.ConceptID,
			Type:         item.Type,
			Code:         item.Code,
			Name:         item.Name,
			Amount:       item.Amount,
			CalculatedAt: item.CalculatedAt.Format(time.RFC3339),
		}
	}

	pr := calc.Payroll
	prResponse := PayrollResponse{
		ID:              pr.ID,
		TenantID:        pr.TenantID,
		EmployeeID:      pr.EmployeeID,
		PeriodStart:     pr.PeriodStart.Format("2006-01-02"),
		PeriodEnd:       pr.PeriodEnd.Format("2006-01-02"),
		PayDate:         pr.PayDate.Format("2006-01-02"),
		GrossAmount:     pr.GrossAmount,
		TotalDeductions: pr.TotalDeductions,
		NetAmount:       pr.NetAmount,
		Status:          pr.Status,
		Items:           items,
		CreatedAt:       pr.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       pr.UpdatedAt.Format(time.RFC3339),
	}

	return &CalculatedPayrollResponse{
		Payroll:         prResponse,
		GrossAmount:     calc.GrossAmount,
		TotalDeductions: calc.TotalDeductions,
		NetAmount:       calc.NetAmount,
		Message:         "Payroll calculated successfully. Review before saving.",
	}
}

// ToPayrollResponse convierte un domain.Payroll a DTO
func ToPayrollResponse(payroll *domain.Payroll, employeeName string) *PayrollResponse {
	items := make([]PayrollItemResponse, len(payroll.Items))
	for i, item := range payroll.Items {
		items[i] = PayrollItemResponse{
			ID:        item.ID,
			ConceptID: item.ConceptID,
			Type:      item.Type,
			Code:      item.Code,
			Name:      item.Name,
			Amount:    item.Amount,
		}
	}

	return &PayrollResponse{
		ID:              payroll.ID,
		TenantID:        payroll.TenantID,
		EmployeeID:      payroll.EmployeeID,
		EmployeeName:    employeeName,
		PeriodStart:     payroll.PeriodStart.Format("2006-01-02"),
		PeriodEnd:       payroll.PeriodEnd.Format("2006-01-02"),
		PayDate:         payroll.PayDate.Format("2006-01-02"),
		GrossAmount:     payroll.GrossAmount,
		TotalDeductions: payroll.TotalDeductions,
		NetAmount:       payroll.NetAmount,
		Status:          payroll.Status,
		Items:           items,
		CreatedAt:       payroll.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       payroll.UpdatedAt.Format(time.RFC3339),
	}
}

// SavePayrollRequest representa el request para guardar una nómina calculada
type SavePayrollRequest struct {
	CalculatePayrollRequest
	Confirm bool `json:"confirm"`
}
