package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
	"github.com/arrase21/crm-users/internal/service"
	"github.com/arrase21/crm-users/internal/transport/http/dto"
	"github.com/arrase21/crm-users/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

type UserHandlerV2 struct {
	svc        *service.UserService
	validator  *middleware.ValidationMiddleware
	tenantFunc func(c *gin.Context) (uint, error)
}

func NewUserHandlerV2(svc *service.UserService) *UserHandlerV2 {
	return &UserHandlerV2{
		svc:        svc,
		validator:  middleware.NewValidationMiddleware(),
		tenantFunc: extractTenantID,
	}
}

// extractTenantID extrae el tenant ID del contexto (implementar según tu auth)
func extractTenantID(c *gin.Context) (uint, error) {
	// Por ahora hardcodeado, pero debería venir del JWT/context
	tenantID := uint(1)

	// En implementación real:
	// tenantID, ok := c.Get("tenant_id")
	// if !ok {
	//     return 0, errors.New("tenant not found")
	// }

	return tenantID, nil
}

// RegisterRoutes registra las rutas con middleware de validación
func (h *UserHandlerV2) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("", h.validator.ValidateBody(&dto.CreateUserRequest{}), h.Create)
		users.GET("", h.List)
		users.GET("/:id", h.validator.ValidateParams(&dto.IDParams{}), h.GetByID)
		users.PUT("/:id", h.validator.ValidateParams(&dto.IDParams{}), h.validator.ValidateBody(&dto.UpdateUserRequest{}), h.Update)
		users.DELETE("/:id", h.validator.ValidateParams(&dto.IDParams{}), h.Delete)
		users.GET("/dni/:dni", h.validator.ValidateParams(&dto.DNIParams{}), h.GetByDni)
	}
}

// Create crea un nuevo usuario con validación robusta
func (h *UserHandlerV2) Create(c *gin.Context) {
	// Extraer tenant ID
	tenantID, err := h.tenantFunc(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:     "unauthorized",
			Message:   "Tenant not found",
			Code:      "TENANT_NOT_FOUND",
			Timestamp: time.Now(),
		})
		return
	}

	// Obtener DTO validado del middleware
	req, exists := c.Get("validated_body")
	if !exists {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:     "validation_error",
			Message:   "Request body not validated",
			Timestamp: time.Now(),
		})
		return
	}

	createReq, ok := req.(*dto.CreateUserRequest)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:     "type_error",
			Message:   "Invalid request type",
			Timestamp: time.Now(),
		})
		return
	}

	// Convertir DTO a modelo de dominio
	user, err := createReq.ToDomain(tenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:     "conversion_error",
			Message:   "Error converting request to domain model: " + err.Error(),
			Timestamp: time.Now(),
		})
		return
	}

	// Validar reglas de negocio adicionales
	if err := user.ValidateAll(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error:     "business_validation_failed",
			Message:   "Business validation failed: " + err.Error(),
			Code:      "BUSINESS_VALIDATION_ERROR",
			Timestamp: time.Now(),
		})
		return
	}

	// Crear usuario
	if err := h.svc.Create(c.Request.Context(), user); err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"user_id": user.ID,
	})
}

// GetByID obtiene un usuario por ID
func (h *UserHandlerV2) GetByID(c *gin.Context) {
	params, _ := c.Get("validated_params")
	idParams := params.(*dto.IDParams)

	user, err := h.svc.GetByID(c.Request.Context(), idParams.ID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// GetByDni obtiene un usuario por DNI
func (h *UserHandlerV2) GetByDni(c *gin.Context) {
	params, _ := c.Get("validated_params")
	dniParams := params.(*dto.DNIParams)

	user, err := h.svc.GetByDni(c.Request.Context(), dniParams.Dni)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// List lista todos los usuarios
func (h *UserHandlerV2) List(c *gin.Context) {
	// Validar query parameters (opcional)
	var queryParams dto.ListUsersQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:     "invalid_query",
			Message:   "Invalid query parameters",
			Timestamp: time.Now(),
		})
		return
	}

	users, err := h.svc.List(c.Request.Context())
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
		"page":  queryParams.Page,
		"limit": queryParams.Limit,
	})
}

// Update actualiza un usuario existente
func (h *UserHandlerV2) Update(c *gin.Context) {
	params, _ := c.Get("validated_params")
	idParams := params.(*dto.IDParams)

	body, _ := c.Get("validated_body")
	updateReq := body.(*dto.UpdateUserRequest)

	// Obtener usuario existente
	user, err := h.svc.GetByID(c.Request.Context(), idParams.ID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	// Aplicar actualizaciones parciales
	if err := updateUserFromDTO(user, updateReq); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:     "update_error",
			Message:   "Error applying updates: " + err.Error(),
			Timestamp: time.Now(),
		})
		return
	}

	// Validar después de las actualizaciones
	if err := user.ValidateAll(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error:     "validation_failed",
			Message:   "Validation failed: " + err.Error(),
			Code:      "VALIDATION_ERROR",
			Timestamp: time.Now(),
		})
		return
	}

	// Actualizar en la base de datos
	if err := h.svc.Update(c.Request.Context(), user); err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
		"user":    user,
	})
}

// Delete elimina un usuario
func (h *UserHandlerV2) Delete(c *gin.Context) {
	params, _ := c.Get("validated_params")
	idParams := params.(*dto.IDParams)

	if err := h.svc.Delete(c.Request.Context(), idParams.ID); err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted successfully",
	})
}

// handleServiceError maneja los errores del servicio de forma estandarizada
func (h *UserHandlerV2) handleServiceError(c *gin.Context, err error) {
	statusCode := http.StatusInternalServerError
	errorCode := "INTERNAL_ERROR"
	message := err.Error()

	switch {
	case err == domain.ErrUserNotFound:
		statusCode = http.StatusNotFound
		errorCode = "USER_NOT_FOUND"
	case err == domain.ErrDniAlreadyExist:
		statusCode = http.StatusConflict
		errorCode = "DNI_ALREADY_EXISTS"
	case err == domain.ErrEmailAlreadyExist:
		statusCode = http.StatusConflict
		errorCode = "EMAIL_ALREADY_EXISTS"
	case err == domain.ErrPhoneAlreadyExist:
		statusCode = http.StatusConflict
		errorCode = "PHONE_ALREADY_EXISTS"
	case err == domain.ErrTenantNotFound:
		statusCode = http.StatusUnauthorized
		errorCode = "TENANT_NOT_FOUND"
	}

	c.JSON(statusCode, dto.ErrorResponse{
		Error:     errorCode,
		Message:   message,
		Code:      errorCode,
		Timestamp: time.Now(),
	})
}

// updateUserFromDTO aplica las actualizaciones del DTO al modelo de dominio
func updateUserFromDTO(user *domain.User, req *dto.UpdateUserRequest) error {
	if req.FirstName != nil {
		user.FirstName = strings.TrimSpace(*req.FirstName)
	}
	if req.LastName != nil {
		user.LastName = strings.TrimSpace(*req.LastName)
	}
	if req.Email != nil {
		user.Email = strings.ToLower(strings.TrimSpace(*req.Email))
	}
	if req.Phone != nil {
		user.Phone = strings.TrimSpace(*req.Phone)
	}
	if req.Dni != nil {
		user.Dni = strings.TrimSpace(*req.Dni)
	}
	if req.Gender != nil {
		user.Gender = strings.ToUpper(strings.TrimSpace(*req.Gender))
	}
	if req.BirthDay != nil {
		birthDate, err := time.Parse("2006-01-02", *req.BirthDay)
		if err != nil {
			return fmt.Errorf("invalid birth_date format: %w", err)
		}
		user.BirthDay = birthDate
	}

	return nil
}
