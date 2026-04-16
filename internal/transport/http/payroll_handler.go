package http

import (
	"net/http"
	"strconv"

	"github.com/arrase21/crm-users/internal/domain"
	"github.com/arrase21/crm-users/internal/service"
	"github.com/arrase21/crm-users/internal/transport/http/dto"
	"github.com/gin-gonic/gin"
)

type PayrollHandler struct {
	calculatorSvc *service.PayrollCalculatorService
	payrollSvc    *service.PayrollService
}

func NewPayrollHandler(calcSvc *service.PayrollCalculatorService, payrollSvc *service.PayrollService) *PayrollHandler {
	return &PayrollHandler{
		calculatorSvc: calcSvc,
		payrollSvc:    payrollSvc,
	}
}

// Calculate calculation una nómina (sin guardar)
// POST /api/v1/payroll/calculate
func (h *PayrollHandler) Calculate(c *gin.Context) {
	var req dto.CalculatePayrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Si no viene pay_date, usar period end
	if req.PayDate.IsZero() {
		req.PayDate = req.PeriodEnd
	}

	serviceReq := service.CalculatePayrollRequest{
		EmployeeID:  req.EmployeeID,
		PeriodStart: req.PeriodStart,
		PeriodEnd:   req.PeriodEnd,
		PayDate:     req.PayDate,
	}

	result, err := h.calculatorSvc.Calculate(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToCalculatedPayrollResponse(result)
	c.JSON(http.StatusOK, response)
}

// CalculateAndSave calcula y guarda una nómina
// POST /api/v1/payroll/calculate-and-save
func (h *PayrollHandler) CalculateAndSave(c *gin.Context) {
	var req dto.SavePayrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Si no viene pay_date, usar period end
	if req.PayDate.IsZero() {
		req.PayDate = req.PeriodEnd
	}

	serviceReq := service.CalculatePayrollRequest{
		EmployeeID:  req.EmployeeID,
		PeriodStart: req.PeriodStart,
		PeriodEnd:   req.PeriodEnd,
		PayDate:     req.PayDate,
	}

	result, err := h.calculatorSvc.CalculateAndSave(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToCalculatedPayrollResponse(result)
	response.Message = "Payroll saved successfully"
	c.JSON(http.StatusOK, response)
}

// GetByID obtiene una nómina por ID
// GET /api/v1/payroll/:id
func (h *PayrollHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	payroll, err := h.payrollSvc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == domain.ErrPayrollNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "payroll not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Obtener nombre del empleado
	employeeName := ""
	if payroll.Employee.User.FirstName != "" {
		employeeName = payroll.Employee.User.FirstName + " " + payroll.Employee.User.LastName
	}

	response := dto.ToPayrollResponse(payroll, employeeName)
	c.JSON(http.StatusOK, response)
}

// ListByEmployee lista las nóminas de un empleado
// GET /api/v1/payroll/employee/:employeeId
func (h *PayrollHandler) ListByEmployee(c *gin.Context) {
	employeeIDStr := c.Param("employeeId")
	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	payrolls, err := h.payrollSvc.ListByEmployee(c.Request.Context(), uint(employeeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if payrolls == nil {
		payrolls = []domain.Payroll{}
	}

	response := make([]dto.PayrollResponse, len(payrolls))
	for i, p := range payrolls {
		employeeName := ""
		if p.Employee.User.FirstName != "" {
			employeeName = p.Employee.User.FirstName + " " + p.Employee.User.LastName
		}
		response[i] = *dto.ToPayrollResponse(&p, employeeName)
	}

	c.JSON(http.StatusOK, gin.H{"payrolls": response})
}

// Delete elimina una nómina
// DELETE /api/v1/payroll/:id
func (h *PayrollHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.payrollSvc.Delete(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payroll deleted"})
}
