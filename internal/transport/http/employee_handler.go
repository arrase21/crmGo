package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/arrase21/crm-users/internal/domain"
	"github.com/arrase21/crm-users/internal/service"
	"github.com/arrase21/crm-users/internal/transport/http/dto"
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
	var req dto.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Sanitize()

	employee := req.ToDomain()

	if err := h.svc.Create(c.Request.Context(), employee); err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
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
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := dto.ToEmployeeResponse(emp)
	c.JSON(http.StatusOK, resp)
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
		if errors.Is(err, domain.ErrActionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := dto.ToEmployeeResponse(emp)
	c.JSON(http.StatusOK, resp)
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
	employeeList := make([]dto.EmployeeResponse, len(employees))
	for i, emp := range employees {
		empResp := dto.ToEmployeeResponse(&emp)
		employeeList[i] = *empResp
	}
	c.JSON(http.StatusOK, gin.H{
		"employees": employeeList,
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

	var req dto.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener empleado existente
	existingEmployee, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Actualizar solo los campos enviados
	if req.DepartmentID != nil {
		existingEmployee.DepartmentID = *req.DepartmentID
	}
	if req.PositionID != nil {
		existingEmployee.PositionID = *req.PositionID
	}
	if req.IsActive != nil {
		existingEmployee.IsActive = *req.IsActive
	}
	if err := h.svc.Update(c.Request.Context(), existingEmployee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "employee updated"})
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not foung"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
