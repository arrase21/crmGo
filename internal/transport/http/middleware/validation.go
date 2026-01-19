package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/arrase21/crm-users/internal/transport/http/dto"
)

// ValidationMiddleware proporciona validación robusta para requests
type ValidationMiddleware struct {
	validator *validator.Validate
}

// NewValidationMiddleware crea una nueva instancia del middleware
func NewValidationMiddleware() *ValidationMiddleware {
	v := validator.New()

	// Registrar validadores personalizados
	v.RegisterValidation("e164", validateE164)
	v.RegisterValidation("dni", validateDNI)
	v.RegisterValidation("alpha_space", validateAlphaSpace)

	return &ValidationMiddleware{validator: v}
}

// ValidateBody valida el body del request contra un struct
func (m *ValidationMiddleware) ValidateBody(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind JSON al struct
		if err := c.ShouldBindJSON(obj); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:     "invalid_json",
				Message:   "Invalid JSON format",
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		// Validar con validator
		if err := m.validator.Struct(obj); err != nil {
			errors := m.formatValidationErrors(err)
			c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Error:       "validation_failed",
				Message:     "Request validation failed",
				Code:        "VALIDATION_ERROR",
				Validations: errors,
				Timestamp:   time.Now(),
			})
			c.Abort()
			return
		}

		// Sanitizar datos si el objeto tiene método Sanitize()
		if sanitizer, ok := obj.(interface{ Sanitize() }); ok {
			sanitizer.Sanitize()
		}

		c.Set("validated_body", obj)
		c.Next()
	}
}

// ValidateQuery valida los parámetros de query
func (m *ValidationMiddleware) ValidateQuery(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindQuery(obj); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:     "invalid_query",
				Message:   "Invalid query parameters",
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		if err := m.validator.Struct(obj); err != nil {
			errors := m.formatValidationErrors(err)
			c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Error:       "query_validation_failed",
				Message:     "Query validation failed",
				Code:        "QUERY_VALIDATION_ERROR",
				Validations: errors,
				Timestamp:   time.Now(),
			})
			c.Abort()
			return
		}

		c.Set("validated_query", obj)
		c.Next()
	}
}

// ValidateParams valida los parámetros de ruta (path params)
func (m *ValidationMiddleware) ValidateParams(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindUri(obj); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:     "invalid_params",
				Message:   "Invalid path parameters",
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		if err := m.validator.Struct(obj); err != nil {
			errors := m.formatValidationErrors(err)
			c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Error:       "params_validation_failed",
				Message:     "Parameters validation failed",
				Code:        "PARAMS_VALIDATION_ERROR",
				Validations: errors,
				Timestamp:   time.Now(),
			})
			c.Abort()
			return
		}

		c.Set("validated_params", obj)
		c.Next()
	}
}

// formatValidationErrors convierte los errores de validator a un formato amigable
func (m *ValidationMiddleware) formatValidationErrors(err error) []dto.ValidationError {
	var errors []dto.ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, dto.ValidationError{
				Field:   e.Field(),
				Message: m.getErrorMessage(e),
				Tag:     e.Tag(),
			})
		}
	}

	return errors
}

// getErrorMessage convierte los tags de validator a mensajes amigables
func (m *ValidationMiddleware) getErrorMessage(fe validator.FieldError) string {
	field := fe.Field()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, fe.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "e164":
		return fmt.Sprintf("%s must be a valid phone number (E.164 format)", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case "numeric":
		return fmt.Sprintf("%s must contain only numbers", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fe.Param())
	case "datetime":
		return fmt.Sprintf("%s must be a valid date in format %s", field, fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// Validadores personalizados

// validateE164 valida números de teléfono en formato E.164
func validateE164(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true // Allow empty for optional fields
	}

	// E.164 format: +[country code][number]
	return strings.HasPrefix(phone, "+") && len(phone) >= 8 && len(phone) <= 15
}

// validateDNI valida formato de DNI argentino
func validateDNI(fl validator.FieldLevel) bool {
	dni := fl.Field().String()
	if len(dni) != 8 {
		return false
	}

	for _, char := range dni {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

// validateAlphaSpace permite letras y espacios
func validateAlphaSpace(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == ' ') {
			return false
		}
	}
	return true
}
