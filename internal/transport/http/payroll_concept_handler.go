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

type PayrollConceptHandler struct {
	svc *service.PayrollConceptService
}

func NewPayrollConceptHandler(svc *service.PayrollConceptService) *PayrollConceptHandler {
	return &PayrollConceptHandler{svc}
}

// Create crea un nuevo concepto de nómina
func (h *PayrollConceptHandler) Create(c *gin.Context) {
	var req dto.CreatePayrollConceptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Sanitize()

	concept := req.ToDomain()

	if err := h.svc.Create(c.Request.Context(), concept); err != nil {
		if errors.Is(err, domain.ErrConceptNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "concept not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "payroll concept created",
		"id":      concept.ID,
	})
}

// GetByID obtiene un concepto por ID
func (h *PayrollConceptHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	concept, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrConceptNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "concept not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dto.ToPayrollConceptResponse(concept)
	c.JSON(http.StatusOK, resp)
}

// GetByCode obtiene un concepto por código
func (h *PayrollConceptHandler) GetByCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	concept, err := h.svc.GetByCode(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, domain.ErrConceptNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "concept not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dto.ToPayrollConceptResponse(concept)
	c.JSON(http.StatusOK, resp)
}

// GetActiveConcepts obtiene todos los conceptos activos
func (h *PayrollConceptHandler) GetActiveConcepts(c *gin.Context) {
	concepts, err := h.svc.GetActiveConcepts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convertir a response
	resp := make([]dto.PayrollConceptResponse, len(concepts))
	for i, concept := range concepts {
		r := dto.ToPayrollConceptResponse(&concept)
		resp[i] = *r
	}

	c.JSON(http.StatusOK, gin.H{"concepts": resp})
}

// List lista todos los conceptos con paginación
func (h *PayrollConceptHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	concepts, total, err := h.svc.List(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	conceptList := make([]dto.PayrollConceptResponse, len(concepts))
	for i, concept := range concepts {
		conceptList[i] = *dto.ToPayrollConceptResponse(&concept)
	}

	c.JSON(http.StatusOK, gin.H{
		"concepts": conceptList,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// Update actualiza un concepto
func (h *PayrollConceptHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	// Obtener concepto existente
	existing, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrConceptNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "concept not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req dto.UpdatePayrollConceptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Actualizar solo los campos enviados
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Type != nil {
		existing.Type = *req.Type
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Percentage != nil {
		existing.Percentage = *req.Percentage
	}
	if req.EmployeePart != nil {
		existing.EmployeePart = *req.EmployeePart
	}
	if req.EmployerPart != nil {
		existing.EmployerPart = *req.EmployerPart
	}
	if req.IsMandatory != nil {
		existing.IsMandatory = *req.IsMandatory
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	if err := h.svc.Update(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payroll concept updated"})
}

// Delete elimina un concepto
func (h *PayrollConceptHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrConceptNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "concept not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// SeedDefaultConcepts carga los conceptos por defecto
func (h *PayrollConceptHandler) SeedDefaultConcepts(c *gin.Context) {
	if err := h.svc.SeedDefaultConcepts(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "default concepts seeded"})
}
