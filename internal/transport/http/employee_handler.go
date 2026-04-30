package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/arrase21/crm-users/internal/domain"
	"github.com/arrase21/crm-users/internal/service"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	svc *service.EmployeeService
}

func NewEmployeeHandler(svc *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{svc}
}

// Create crea un nuevo empleado
func (h *EmployeeHandler) Create(c *gin.Context) {
	var req service.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if errors.Is(err, service.ErrDuplicateEmployee) {
			c.JSON(http.StatusConflict, gin.H{"error": "user already has an employee"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "employee created",
		"id":      employee.ID,
	})
}

// GetByID obtiene un empleado por ID
func (h *EmployeeHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 31)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	emp, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, emp)
}

func (h *EmployeeHandler) GetByUserID(c *gin.Context) {
	userIDstr := c.Query("user_id")
	userID, err := strconv.ParseUint(userIDstr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	emp, err := h.svc.GetByUserID(c.Request.Context(), uint(userID))
	if err != nil {
		if errors.Is(err, service.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, emp)
}

func (h *EmployeeHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	employees, total, err := h.svc.List(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"employees": employees,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req service.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEmployee, err := h.svc.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		if errors.Is(err, service.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedEmployee)
}

// ChangeStatus cambia el estado del empleado
func (h *EmployeeHandler) ChangeStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req service.ChangeEmployeeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEmployee, err := h.svc.ChangeStatus(c.Request.Context(), uint(id), req)
	if err != nil {
		if errors.Is(err, service.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedEmployee)
}

// GetStatistics obtiene estadísticas de empleados
func (h *EmployeeHandler) GetStatistics(c *gin.Context) {
	stats, err := h.svc.GetStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	// Usar Deactivate en vez de Delete para validaciones
	if err := h.svc.Deactivate(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, service.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}
		if errors.Is(err, service.ErrCannotDeactivateEmployee) {
			c.JSON(http.StatusConflict, gin.H{"error": "cannot deactivate: has active contract"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employee deactivated"})
}

// Helper para convertir domain.Employee a response (para el caso GetByUserID que aún usa el domain)
func toEmployeeResponseFromDomain(emp *domain.Employee) gin.H {
	resp := gin.H{
		"id":        emp.ID,
		"is_active": emp.IsActive,
		"user": gin.H{
			"id":         emp.User.ID,
			"first_name": emp.User.FirstName,
			"last_name":  emp.User.LastName,
			"dni":        emp.User.Dni,
			"email":      emp.User.Email,
		},
	}

	if emp.DepartmentID != 0 {
		resp["department"] = gin.H{"id": emp.Department.ID, "name": emp.Department.Name}
	}
	if emp.PositionID != 0 {
		resp["position"] = gin.H{"id": emp.Position.ID, "name": emp.Position.NamePosition}
	}

	return resp
}
