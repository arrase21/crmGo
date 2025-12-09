package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
	"github.com/arrase21/crm-users/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc}
}

type CreateUserRequest struct {
	FirstName string `json:"firstname" binding:"required,max=30"`
	LastName  string `json:"lastname" binding:"required,max=40"`
	Dni       string `json:"dni" binding:"required,max=20"`
	Gender    string `json:"gender" binding:"required,oneof=M F"`
	Phone     string `json:"phone" binding:"required,max=15"`
	Email     string `json:"email" binding:"required,email,max=50"`
	BirthDay  string `json:"birth_day" binding:"required"`
}

func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parsear la fecha correctamente
	birth, err := time.Parse("2006-01-02", req.BirthDay)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid birth_day format (expected YYYY-MM-DD)"})
		return
	}

	// Construir objeto de dominio
	user := &domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Dni:       req.Dni,
		Gender:    req.Gender,
		Phone:     req.Phone,
		Email:     req.Email,
		BirthDay:  birth,
	}

	user.Normalize()

	// Validaciones del dominio
	if err := user.ValidateAll(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Guardar en el servicio
	if err := h.svc.Create(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	usr, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, usr)
}

func (h *UserHandler) GetByDni(c *gin.Context) {
	dni := c.Query("value")
	if dni == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dni is required"})
		return
	}
	usr, err := h.svc.GetByDni(c.Request.Context(), dni)
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, usr)
}

func (h *UserHandler) List(c *gin.Context) {
	usrs, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": usrs})
}

type UpdateUserRequest struct {
	FirstName *string `json:"firstname"`
	LastName  *string `json:"lastname"`
	Dni       *string `json:"dni"`
	Gender    *string `json:"gender"`
	Phone     *string `json:"phone"`
	Email     *string `json:"email"`
	BirthDay  *string `json:"birth_day"`
}

func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
}
