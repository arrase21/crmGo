# VALIDACIÓN ROBUSTA - Guía de Implementación

## 🎯 ¿QUÉ IMPLEMENTAMOS?

Creamos un sistema de validación **multi-capa** que protege tu CRM en múltiples niveles:

### **1. DTO Layer (Transport)**
```json
{
  "first_name": "John",
  "email": "john@example.com",
  "dni": "12345678",
  "birth_day": "1990-01-01"
}
```

**Validaciones automáticas:**
- `required` - Campo obligatorio
- `email` - Formato de email válido
- `len=8` - Longitud exacta de 8 caracteres
- `numeric` - Solo números
- `datetime=2006-01-02` - Formato de fecha

### **2. Middleware Layer**
Validación antes de llegar al handler:
```go
//middleware.ValidationMiddleware
func (m *ValidationMiddleware) ValidateBody(obj interface{}) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Bind JSON
        if err := c.ShouldBindJSON(obj); err != nil {
            c.JSON(400, ErrorResponse{Error: "invalid_json"})
            return
        }
        
        // 2. Validar con validator tags
        if err := m.validator.Struct(obj); err != nil {
            errors := m.formatValidationErrors(err)
            c.JSON(422, ErrorResponse{
                Error: "validation_failed",
                Validations: errors
            })
            return
        }
        
        // 3. Sanitizar datos
        if sanitizer, ok := obj.(interface{ Sanitize() }); ok {
            sanitizer.Sanitize()
        }
        
        c.Set("validated_body", obj)
        c.Next()
    }
}
```

### **3. Service Layer**
Validación de lógica de negocio:
```go
func (s *UserService) Create(ctx context.Context, req *CreateUserRequest) error {
    // 1. Validar reglas de negocio
    if !isValidDNI(req.Dni) {
        return ErrInvalidDNIFormat
    }
    
    // 2. Validar unicidad
    if exists, _ := s.repo.GetByDni(ctx, req.Dni); exists != nil {
        return ErrDNIAlreadyExists
    }
    
    // 3. Validar edad
    if calculateAge(req.BirthDay) < 18 {
        return ErrUserMustBeAdult
    }
}
```

### **4. Domain Layer**
Validación de consistencia:
```go
func (u *User) Validate() error {
    if u.Gender != "M" && u.Gender != "F" {
        return ErrInvalidGender
    }
    
    if u.IsMinor() {
        return ErrUserMustBeAdult
    }
    
    return nil
}
```

## 🛡️ EJEMPLOS DE VALIDACIÓN EN ACCIÓN

### **Caso 1: Request válido**
```bash
POST /api/v1/users
Content-Type: application/json

{
  "first_name": "Juan",
  "last_name": "Pérez",
  "dni": "12345678",
  "gender": "M",
  "phone": "+5491123456789",
  "email": "juan.perez@example.com",
  "birth_day": "1990-01-01"
}
```

**Response:**
```json
{
  "message": "user created successfully",
  "user_id": 1
}
```

### **Caso 2: Datos inválidos**
```bash
POST /api/v1/users
Content-Type: application/json

{
  "first_name": "J",
  "email": "invalid-email",
  "dni": "123",
  "gender": "X",
  "phone": "123456",
  "birth_day": "2025-01-01"
}
```

**Response:**
```json
{
  "error": "validation_failed",
  "message": "Request validation failed",
  "code": "VALIDATION_ERROR",
  "validations": [
    {
      "field": "FirstName",
      "message": "FirstName must be at least 2 characters",
      "tag": "min"
    },
    {
      "field": "Email",
      "message": "Email must be a valid email address",
      "tag": "email"
    },
    {
      "field": "Dni",
      "message": "Dni must be exactly 8 characters",
      "tag": "len"
    },
    {
      "field": "Gender",
      "message": "Gender must be one of: M F",
      "tag": "oneof"
    },
    {
      "field": "Phone",
      "message": "Phone must be a valid phone number (E.164 format)",
      "tag": "e164"
    },
    {
      "field": "BirthDay",
      "message": "BirthDay must be a valid date in format 2006-01-02",
      "tag": "datetime"
    }
  ],
  "timestamp": "2024-01-19T15:30:00Z"
}
```

### **Caso 3: Error de negocio**
```bash
POST /api/v1/users
Content-Type: application/json

{
  "first_name": "Juan",
  "last_name": "Pérez",
  "dni": "12345678",
  "gender": "M",
  "phone": "+5491123456789",
  "email": "juan.perez@example.com",
  "birth_day": "2010-01-01"  // Menor de 18 años
}
```

**Response:**
```json
{
  "error": "business_validation_failed",
  "message": "Business validation failed: user must be over 18",
  "code": "BUSINESS_VALIDATION_ERROR",
  "timestamp": "2024-01-19T15:30:00Z"
}
```

## 🔧 VALIDADORES PERSONALIZADOS

### **Validador de DNI Argentino**
```go
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
```

### **Validador de Teléfono E.164**
```go
func validateE164(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    
    // E.164 format: +[country code][number]
    return strings.HasPrefix(phone, "+") && len(phone) >= 8 && len(phone) <= 15
}
```

## 🎯 BENEFICIOS DE LA VALIDACIÓN ROBUSTA

### **1. Seguridad**
- **Input sanitization**: Limpieza automática de datos
- **Injection prevention**: Protección contra SQLi, XSS
- **Data integrity**: Garantiza consistencia

### **2. Experiencia de Usuario**
- **Feedback inmediato**: Errores descriptivos
- **Validación temprana**: Antes del procesamiento
- **Mensajes útiles**: Ayuda al usuario a corregir

### **3. Mantenibilidad**
- **Centralización**: Lógica de validación en un solo lugar
- **Reutilizabilidad**: DTOs reutilizables
- **Testabilidad**: Cada capa es testeable

### **4. Performance**
- **Early returns**: Falla rápido
- **Minimal processing**: No procesa datos inválidos
- **Efficient validation**: Validaciones optimizadas

## 🚀 CÓMO USARLO

### **1. Crear DTO con validation tags**
```go
type CreateUserRequest struct {
    FirstName string `json:"first_name" validate:"required,min=2,max=30,alpha"`
    Email     string `json:"email" validate:"required,email"`
    Dni       string `json:"dni" validate:"required,len=8,numeric,dni"`
}
```

### **2. Aplicar middleware en rutas**
```go
func (h *UserHandlerV2) RegisterRoutes(router *gin.RouterGroup) {
    users := router.Group("/users")
    {
        users.POST("", 
            h.validator.ValidateBody(&dto.CreateUserRequest{}), 
            h.Create)
        users.PUT("/:id", 
            h.validator.ValidateParams(&dto.IDParams{}),
            h.validator.ValidateBody(&dto.UpdateUserRequest{}), 
            h.Update)
    }
}
```

### **3. Implementar Sanitize() en DTO**
```go
func (r *CreateUserRequest) Sanitize() {
    r.FirstName = strings.TrimSpace(r.FirstName)
    r.Email = strings.ToLower(strings.TrimSpace(r.Email))
    r.Dni = strings.TrimSpace(r.Dni)
    r.Gender = strings.ToUpper(r.Gender)
}
```

## 📊 COMPARATIVA ANTES vs DESPUÉS

| Aspecto | ❌ Antes | ✅ Después |
|---------|----------|------------|
| Validación | Manual en handler | Automática multi-capa |
| Errores | Genéricos | Específicos y útiles |
| Sanitización | Manual | Automática |
| Testing | Difícil | Fácil con mocks |
| Mantenimiento | Spaghetti | Centralizado |
| Seguridad | Baja | Alta |
| UX | Pobre | Excelente |

## 🎯 PRÓXIMOS PASOS

1. **Agregar más validadores personalizados**
   - CUIL/CUIT argentino
   - Validación de códigos postales
   - Validación de provincias

2. **Implementar rate limiting por tenant**
3. **Agregar logging estructurado**
4. **Crear tests de integración E2E**
5. **Implementar OpenAPI documentation**

¡Tu CRM ahora tiene **validación de nivel enterprise**! 🚀