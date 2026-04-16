package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/arrase21/crm-users/internal/service"
	"github.com/gin-gonic/gin"
)

// PayrollStateHandler maneja las transiciones de estado y batch processing
type PayrollStateHandler struct {
	stateSvc   *service.PayrollStateService
	batchSvc   *service.PayrollBatchService
	payrollSvc *service.PayrollService
}

func NewPayrollStateHandler(
	stateSvc *service.PayrollStateService,
	batchSvc *service.PayrollBatchService,
	payrollSvc *service.PayrollService,
) *PayrollStateHandler {
	return &PayrollStateHandler{
		stateSvc:   stateSvc,
		batchSvc:   batchSvc,
		payrollSvc: payrollSvc,
	}
}

// MarkAsPaidRequest representa el request para marcar como pagada
type MarkAsPaidRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required"`
}

// MarkAsPaid marca una nómina como pagada
// POST /api/v1/payroll/:id/mark-paid
func (h *PayrollStateHandler) MarkAsPaid(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payroll id"})
		return
	}

	var req MarkAsPaidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment_method is required"})
		return
	}

	payment, err := h.stateSvc.MarkAsPaid(c.Request.Context(), uint(id), req.PaymentMethod)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrInvalidStatusTransition ||
			err == service.ErrPayrollNotInDraft {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "payroll marked as paid successfully",
		"payment_id":     payment.ID,
		"payment_method": payment.Method,
		"amount":         payment.Amount,
		"paid_at":        payment.PaidAt.Format("2006-01-02 15:04:05"),
	})
}

// RevertToDraft revierte una nómina a estado draft
// POST /api/v1/payroll/:id/revert-to-draft
func (h *PayrollStateHandler) RevertToDraft(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payroll id"})
		return
	}

	err = h.stateSvc.RevertToDraft(c.Request.Context(), uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrInvalidStatusTransition {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payroll reverted to draft"})
}

// GetPaymentInfo obtiene la información de pago de una nómina
// GET /api/v1/payroll/:id/payment
func (h *PayrollStateHandler) GetPaymentInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payroll id"})
		return
	}

	payment, err := h.stateSvc.GetPaymentHistory(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no payment record found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment": gin.H{
			"id":             payment.ID,
			"method":         payment.Method,
			"bank_name":      payment.BankName,
			"account_number": payment.AccountNumber,
			"amount":         payment.Amount,
			"paid_at":        payment.PaidAt.Format("2006-01-02 15:04:05"),
			"status":         payment.Status,
		},
	})
}

// BatchPayrollRequest representa el request para batch processing
type BatchPayrollRequest struct {
	PeriodStart string `json:"period_start" binding:"required"`
	PeriodEnd   string `json:"period_end" binding:"required"`
	PayDate     string `json:"pay_date"`
	Concurrent  bool   `json:"concurrent"`
	Departments []uint `json:"department_ids"`
}

// ProcessBatch genera nóminas para todos los empleados activos
// POST /api/v1/payroll/batch
func (h *PayrollStateHandler) ProcessBatch(c *gin.Context) {
	var req BatchPayrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	periodStart, err := parseDate(req.PeriodStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_start format, use YYYY-MM-DD"})
		return
	}

	periodEnd, err := parseDate(req.PeriodEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_end format, use YYYY-MM-DD"})
		return
	}

	var payDate time.Time
	if req.PayDate != "" {
		payDate, err = parseDate(req.PayDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pay_date format, use YYYY-MM-DD"})
			return
		}
	}

	batchReq := service.BatchPayrollRequest{
		PeriodStart:   periodStart,
		PeriodEnd:     periodEnd,
		PayDate:       payDate,
		Concurrent:    req.Concurrent,
		DepartmentIDs: req.Departments,
	}

	result, err := h.batchSvc.CalculatePeriodSummary(c.Request.Context(), batchReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetPayrollSummary obtiene un resumen de todas las nóminas de un periodo
// GET /api/v1/payroll/summary?period_start=2024-01-01&period_end=2024-01-31
func (h *PayrollStateHandler) GetPayrollSummary(c *gin.Context) {
	periodStartStr := c.Query("period_start")
	periodEndStr := c.Query("period_end")

	if periodStartStr == "" || periodEndStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period_start and period_end are required"})
		return
	}

	periodStart, err := parseDate(periodStartStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_start format"})
		return
	}

	periodEnd, err := parseDate(periodEndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_end format"})
		return
	}

	payrolls, err := h.payrollSvc.GetByPeriod(c.Request.Context(), periodStart, periodEnd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calcular resumen
	var summary struct {
		TotalCount      int     `json:"total_count"`
		DraftCount      int     `json:"draft_count"`
		CalculatedCount int     `json:"calculated_count"`
		PaidCount       int     `json:"paid_count"`
		TotalGross      float64 `json:"total_gross_amount"`
		TotalDeductions float64 `json:"total_deductions"`
		TotalNet        float64 `json:"total_net_amount"`
	}

	for _, p := range payrolls {
		summary.TotalCount++
		switch p.Status {
		case "draft":
			summary.DraftCount++
		case "calculated":
			summary.CalculatedCount++
		case "paid":
			summary.PaidCount++
		}
		summary.TotalGross += p.GrossAmount
		summary.TotalDeductions += p.TotalDeductions
		summary.TotalNet += p.NetAmount
	}

	c.JSON(http.StatusOK, gin.H{
		"period_start": periodStartStr,
		"period_end":   periodEndStr,
		"summary":      summary,
	})
}

func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}
