package service

import (
	"context"
	"fmt"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
)

// ========================================
// Report Service - Generación de reportes de nómina
// ========================================

type PayrollReportService struct {
	payrollRepo     domain.PayrollRepo
	employeeRepo    domain.EmployeeRepo
	contractRepo    domain.EmployeeContractRepo
	paymentRepo     domain.PaymentRepo
	payrollItemRepo domain.PayrollItemRepo
}

func NewPayrollReportService(
	payrollRepo domain.PayrollRepo,
	employeeRepo domain.EmployeeRepo,
	contractRepo domain.EmployeeContractRepo,
	paymentRepo domain.PaymentRepo,
	payrollItemRepo domain.PayrollItemRepo,
) *PayrollReportService {
	return &PayrollReportService{
		payrollRepo:     payrollRepo,
		employeeRepo:    employeeRepo,
		contractRepo:    contractRepo,
		paymentRepo:     paymentRepo,
		payrollItemRepo: payrollItemRepo,
	}
}

// ========================================
// Tipos de Reporte
// ========================================

// PayrollBookReport representa el libro de nóminas
type PayrollBookReport struct {
	PeriodStart     string              `json:"period_start"`
	PeriodEnd       string              `json:"period_end"`
	GeneratedAt     string              `json:"generated_at"`
	TotalEmployees  int                 `json:"total_employees"`
	TotalGross      float64             `json:"total_gross"`
	TotalDeductions float64             `json:"total_deductions"`
	TotalNet        float64             `json:"total_net"`
	Employees       []PayrollBookDetail `json:"employees"`
}

type PayrollBookDetail struct {
	EmployeeID   uint    `json:"employee_id"`
	EmployeeName string  `json:"employee_name"`
	Document     string  `json:"document"`
	BaseSalary   float64 `json:"base_salary"`
	Overtime     float64 `json:"overtime"`
	Bonuses      float64 `json:"bonuses"`
	GrossAmount  float64 `json:"gross_amount"`
	Deductions   float64 `json:"deductions"`
	NetAmount    float64 `json:"net_amount"`
}

// EmployeePayrollReport representa el historial de nómina de un empleado
type EmployeePayrollReport struct {
	EmployeeID    uint             `json:"employee_id"`
	EmployeeName  string           `json:"employee_name"`
	Document      string           `json:"document"`
	PeriodStart   string           `json:"period_start"`
	PeriodEnd     string           `json:"period_end"`
	TotalPayrolls int              `json:"total_payrolls"`
	TotalEarned   float64          `json:"total_earned"`
	TotalDeducted float64          `json:"total_deducted"`
	TotalNet      float64          `json:"total_net"`
	Payrolls      []PayrollSummary `json:"payrolls"`
}

type PayrollSummary struct {
	Period      string  `json:"period"`
	PayDate     string  `json:"pay_date"`
	GrossAmount float64 `json:"gross_amount"`
	Deductions  float64 `json:"deductions"`
	NetAmount   float64 `json:"net_amount"`
	Status      string  `json:"status"`
}

// IncomeCertification representa una certificación de ingresos
type IncomeCertification struct {
	EmployeeID      uint    `json:"employee_id"`
	EmployeeName    string  `json:"employee_name"`
	Document        string  `json:"document"`
	PeriodStart     string  `json:"period_start"`
	PeriodEnd       string  `json:"period_end"`
	TotalIncome     float64 `json:"total_income"`
	TotalDeductions float64 `json:"total_deductions"`
	TotalNet        float64 `json:"total_net"`
	AvgMonthly      float64 `json:"avg_monthly_income"`
	IssuedAt        string  `json:"issued_at"`
}

// DeductionsReport representa el reporte de deducciones
type DeductionsReport struct {
	PeriodStart  string            `json:"period_start"`
	PeriodEnd    string            `json:"period_end"`
	GeneratedAt  string            `json:"generated_at"`
	TotalHealth  float64           `json:"total_health"`
	TotalPension float64           `json:"total_pension"`
	TotalTax     float64           `json:"total_tax"`
	TotalOther   float64           `json:"total_other"`
	ByEmployee   []DeductionDetail `json:"by_employee"`
}

type DeductionDetail struct {
	EmployeeID   uint    `json:"employee_id"`
	EmployeeName string  `json:"employee_name"`
	Health       float64 `json:"health"`
	Pension      float64 `json:"pension"`
	Tax          float64 `json:"tax"`
	Other        float64 `json:"other"`
	Total        float64 `json:"total"`
}

// ========================================
// Métodos de Reporte
// ========================================

// GeneratePayrollBook genera el libro de nóminas para un período
func (s *PayrollReportService) GeneratePayrollBook(ctx context.Context, periodStart, periodEnd time.Time) (*PayrollBookReport, error) {
	payrolls, err := s.payrollRepo.GetByPeriod(ctx, periodStart, periodEnd)
	if err != nil {
		return nil, err
	}

	report := &PayrollBookReport{
		PeriodStart: periodStart.Format("2006-01-02"),
		PeriodEnd:   periodEnd.Format("2006-01-02"),
		GeneratedAt: time.Now().Format(time.RFC3339),
		Employees:   make([]PayrollBookDetail, 0, len(payrolls)),
	}

	for _, p := range payrolls {
		detail := PayrollBookDetail{
			EmployeeID:   p.EmployeeID,
			EmployeeName: p.Employee.User.FirstName + " " + p.Employee.User.LastName,
			Document:     p.Employee.User.Dni,
			GrossAmount:  p.GrossAmount,
			Deductions:   p.TotalDeductions,
			NetAmount:    p.NetAmount,
		}
		report.Employees = append(report.Employees, detail)
		report.TotalGross += p.GrossAmount
		report.TotalDeductions += p.TotalDeductions
		report.TotalNet += p.NetAmount
	}

	report.TotalEmployees = len(payrolls)
	return report, nil
}

// GenerateEmployeePayrollReport genera historial de nómina de un empleado
func (s *PayrollReportService) GenerateEmployeePayrollReport(ctx context.Context, employeeID uint, periodStart, periodEnd time.Time) (*EmployeePayrollReport, error) {
	payrolls, err := s.payrollRepo.ListByEmployee(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	if len(payrolls) == 0 {
		return nil, fmt.Errorf("no payrolls found for employee %d", employeeID)
	}

	// Filtrar por período
	var filteredPayrolls []domain.Payroll
	for _, p := range payrolls {
		if !p.PeriodStart.Before(periodStart) && !p.PeriodEnd.After(periodEnd) {
			filteredPayrolls = append(filteredPayrolls, p)
		}
	}

	emp := payrolls[0].Employee
	report := &EmployeePayrollReport{
		EmployeeID:   employeeID,
		EmployeeName: emp.User.FirstName + " " + emp.User.LastName,
		Document:     emp.User.Dni,
		PeriodStart:  periodStart.Format("2006-01-02"),
		PeriodEnd:    periodEnd.Format("2006-01-02"),
		Payrolls:     make([]PayrollSummary, 0, len(filteredPayrolls)),
	}

	for _, p := range filteredPayrolls {
		summary := PayrollSummary{
			Period:      fmt.Sprintf("%s - %s", p.PeriodStart.Format("2006-01-02"), p.PeriodEnd.Format("2006-01-02")),
			PayDate:     p.PayDate.Format("2006-01-02"),
			GrossAmount: p.GrossAmount,
			Deductions:  p.TotalDeductions,
			NetAmount:   p.NetAmount,
			Status:      p.Status,
		}
		report.Payrolls = append(report.Payrolls, summary)
		report.TotalEarned += p.GrossAmount
		report.TotalDeducted += p.TotalDeductions
		report.TotalNet += p.NetAmount
	}

	report.TotalPayrolls = len(filteredPayrolls)
	return report, nil
}

// GenerateIncomeCertification genera certificación de ingresos
func (s *PayrollReportService) GenerateIncomeCertification(ctx context.Context, employeeID uint, year int) (*IncomeCertification, error) {
	periodStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	payrolls, err := s.payrollRepo.ListByEmployee(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	if len(payrolls) == 0 {
		return nil, fmt.Errorf("no payrolls found for employee %d in year %d", employeeID, year)
	}

	var totalIncome, totalDeductions, totalNet float64
	var emp domain.Employee

	for _, p := range payrolls {
		if !p.PeriodStart.Before(periodStart) && !p.PeriodEnd.After(periodEnd) {
			totalIncome += p.GrossAmount
			totalDeductions += p.TotalDeductions
			totalNet += p.NetAmount
			emp = p.Employee
		}
	}

	months := 12
	if year == time.Now().Year() {
		months = int(time.Now().Month())
	}

	return &IncomeCertification{
		EmployeeID:      employeeID,
		EmployeeName:    emp.User.FirstName + " " + emp.User.LastName,
		Document:        emp.User.Dni,
		PeriodStart:     periodStart.Format("2006-01-02"),
		PeriodEnd:       periodEnd.Format("2006-01-02"),
		TotalIncome:     totalIncome,
		TotalDeductions: totalDeductions,
		TotalNet:        totalNet,
		AvgMonthly:      totalIncome / float64(months),
		IssuedAt:        time.Now().Format(time.RFC3339),
	}, nil
}

// GenerateDeductionsReport genera reporte de deducciones
func (s *PayrollReportService) GenerateDeductionsReport(ctx context.Context, periodStart, periodEnd time.Time) (*DeductionsReport, error) {
	payrolls, err := s.payrollRepo.GetByPeriod(ctx, periodStart, periodEnd)
	if err != nil {
		return nil, err
	}

	report := &DeductionsReport{
		PeriodStart: periodStart.Format("2006-01-02"),
		PeriodEnd:   periodEnd.Format("2006-01-02"),
		GeneratedAt: time.Now().Format(time.RFC3339),
		ByEmployee:  make([]DeductionDetail, 0),
	}

	for _, p := range payrolls {
		// Obtener los items de nómina para identificar deducciones
		items, _ := s.payrollItemRepo.GetByIDPayrollID(ctx, p.ID)

		var health, pension, tax, other float64
		for _, item := range items {
			switch item.Code {
			case domain.ConceptHealth:
				health = item.Amount
			case domain.ConceptPension:
				pension = item.Amount
			case domain.ConceptTax:
				tax = item.Amount
			default:
				if item.Type == domain.PayrollTypeDeduction {
					other += item.Amount
				}
			}
		}

		detail := DeductionDetail{
			EmployeeID:   p.EmployeeID,
			EmployeeName: p.Employee.User.FirstName + " " + p.Employee.User.LastName,
			Health:       health,
			Pension:      pension,
			Tax:          tax,
			Other:        other,
			Total:        p.TotalDeductions,
		}

		report.ByEmployee = append(report.ByEmployee, detail)
		report.TotalHealth += health
		report.TotalPension += pension
		report.TotalTax += tax
		report.TotalOther += other
	}

	return report, nil
}
