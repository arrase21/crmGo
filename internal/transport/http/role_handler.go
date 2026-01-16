package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/arrase21/crm-users/internal/domain"
	"github.com/arrase21/crm-users/internal/service"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleSvc       *service.RoleService
	permissionSvc *service.PermissionService
}

func NewRoleHandler(roleSvc *service.RoleService, permissionSvc *service.PermissionService) *RoleHandler {
	return &RoleHandler{
		roleSvc:       roleSvc,
		permissionSvc: permissionSvc,
	}
}

// ========================================
// CRUD de Roles
// ========================================

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description" binding:"max=255"`
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := &domain.Role{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}

	if err := h.roleSvc.Create(c.Request.Context(), role); err != nil {
		if errors.Is(err, domain.ErrRoleExisting) {
			c.JSON(http.StatusConflict, gin.H{"error": "role already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "role created",
		"role":    role,
	})
}

func (h *RoleHandler) List(c *gin.Context) {
	roles, err := h.roleSvc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

func (h *RoleHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	role, err := h.roleSvc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrRoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

type UpdateRoleRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=50"`
	Description *string `json:"description" binding:"omitempty,max=255"`
	IsActive    *bool   `json:"is_active"`
}

func (h *RoleHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.roleSvc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrRoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Actualizar campos
	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}
	if req.IsActive != nil {
		role.IsActive = *req.IsActive
	}

	if err := h.roleSvc.Update(c.Request.Context(), role); err != nil {
		if errors.Is(err, domain.ErrRoleExisting) {
			c.JSON(http.StatusConflict, gin.H{"error": "role name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "role updated",
		"role":    role,
	})
}

func (h *RoleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.roleSvc.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrRoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ========================================
// Gestión de permisos de roles
// ========================================

type AssignPermissionRequest struct {
	ActionID uint `json:"action_id" binding:"required"`
}

func (h *RoleHandler) AssignPermission(c *gin.Context) {
	roleIDStr := c.Param("id") // ✅ Cambiado de "id" a "roleId"
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	var req AssignPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.permissionSvc.AssignPermissionToRole(c.Request.Context(), uint(roleID), req.ActionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission assigned"})
}

func (h *RoleHandler) RevokePermission(c *gin.Context) {
	roleIDStr := c.Param("id") // ✅ Cambiado de "id" a "roleId"
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	actionIDStr := c.Param("actionId")
	actionID, err := strconv.ParseUint(actionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action id"})
		return
	}

	if err := h.permissionSvc.RevokePermissionFromRole(c.Request.Context(), uint(roleID), uint(actionID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission revoked"})
}

// ========================================
// Asignación de roles a usuarios
// ========================================

type AssignRoleToUserRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	RoleID uint `json:"role_id" binding:"required"`
}

func (h *RoleHandler) AssignRoleToUser(c *gin.Context) {
	var req AssignRoleToUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.permissionSvc.AssignRoleToUser(c.Request.Context(), req.UserID, req.RoleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role assigned to user"})
}

type RevokeRoleFromUserRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	RoleID uint `json:"role_id" binding:"required"`
}

func (h *RoleHandler) RevokeRoleFromUser(c *gin.Context) {
	var req RevokeRoleFromUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.permissionSvc.RevokeRoleFromUser(c.Request.Context(), req.UserID, req.RoleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role revoked from user"})
}
