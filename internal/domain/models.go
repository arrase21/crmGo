package domain

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// ========================================
// Payroll Constants
// ========================================

const (
	PayrollStatusDraft      = "draft"
	PayrollStatusCalculated = "calculated"
	PayrollStatusPaid       = "paid"
	PayrollStatusApproved   = "approved"
)

// Employee status constants
const (
	EmployeeStatusActive    = "active"
	EmployeeStatusInactive  = "inactive"
	EmployeeStatusSuspended = "suspended"
)

const (
	PayrollTypeEarning              = "earning"
	PayrollTypeDeduction            = "deduction"
	PayrollTypeEmployerContribution = "employer_contribution"
)

// PayrollConcept codes
const (
	ConceptBaseSalary      = "BASE_SALARY"
	ConceptTransport       = "TRANSPORT"
	ConceptHousing         = "HOUSING"
	ConceptOvertime        = "OVERTIME"
	ConceptBonus           = "BONUS"
	ConceptHealth          = "HEALTH"
	ConceptPension         = "PENSION"
	ConceptTax             = "TAX"
	ConceptOtherDeduction  = "OTHER_DEDUCTION"
	ConceptHealthEmployer  = "HEALTH_EMPLOYER"
	ConceptPensionEmployer = "PENSION_EMPLOYER"
	ConceptParafiscales    = "PARAFISCALES"
)

type User struct {
	ID uint `gorm:"primaryKey" json:"id"`
	// TenantID uint `gorm:"not null;uniqueIndex:idx_users_tenant_dni;uniqueIndex:idx_users_tenant_phone;uniqueIndex:idx_users_tenant_email"`
	TenantID  uint           `gorm:"not null;uniqueIndex:idx_users_tenant_id" json:"tenant_id"`
	FirstName string         `gorm:"size:30;not null" json:"first_name"`
	LastName  string         `gorm:"size:40;not null" json:"last_name"`
	Dni       string         `gorm:"size:20;not null;uniqueIndex:idx_users_tenant_dni,composite:tenant_dni" json:"dni"`
	Gender    string         `gorm:"size:1;not null;check:gender IN ('M', 'F')" json:"gender"`
	Phone     string         `gorm:"size:15;not null;uniqueIndex:idx_users_tenant_phone,composite:tenant_phone" json:"phone"`
	Email     string         `gorm:"size:50;not null;uniqueIndex:idx_users_tenant_email,composite:tenant_email" json:"email"`
	BirthDay  time.Time      `gorm:"not null" json:"birth_day"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_users_deleted_at" json:"deleted_at,omitzero"`
	Roles     []Role         `gorm:"many2many:user_roles;" json:"roles,omitzero"`
}

// Permission
type Permission struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	Name        string             `gorm:"size:50;not null;unique" json:"name"`
	DisplayName string             `gorm:"size:100" json:"display_name"`
	Description string             `gorm:"size:255" json:"description"`
	Module      string             `gorm:"size:50" json:"module"`
	IsActive    bool               `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time          `gorm:"autoCreateTime" json:"created_at"`
	Actions     []PermissionAction `gorm:"foreignKey:ResourceID;references:ID" json:"actions,omitzero"`
}

// Actions
type PermissionAction struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ResourceID  uint       `gorm:"not null;index" json:"resource_id"`
	Action      string     `gorm:"size:20;not null" json:"action"`
	DisplayName string     `gorm:"size:100" json:"display_name"`
	Description string     `gorm:"size:255" json:"description"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	Resource    Permission `gorm:"foreignKey:ResourceID" json:"resource,omitzero"`
}

// Roles
type Role struct {
	ID              uint             `gorm:"primaryKey" json:"id"`
	TenantID        uint             `gorm:"not null;uniqueIndex:index_role_tenant_name" json:"tenant_id"`
	Name            string           `gorm:"size:50;not null;uniqueIndex:idx_role_tenant_name,composite:tenant_name" json:"name"`
	Description     string           `gorm:"size:255" json:"description"`
	IsSystem        bool             `gorm:"default:false" json:"is_system"`
	IsActive        bool             `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt   `gorm:"index" json:"deleted_at,omitzero"`
	RolePermissions []RolePermission `gorm:"foreignKey:RoleID" json:"permissions,omitempty"`
	Users           []User           `gorm:"many2many:user_roles;" json:"-"`
}

type RolePermission struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	RoleID    uint             `gorm:"not null;uniqueIndex:idx_role_action" json:"role_id"`
	ActionID  uint             `gorm:"not null;uniqueIndex:idx_role_action" json:"action_id"`
	GrantedAt time.Time        `gorm:"autoCreateTime" json:"granted_at"`
	Role      Role             `gorm:"foreignKey:RoleID" json:"-"`
	Action    PermissionAction `gorm:"foreignKey:ActionID" json:"action,omitzero"`
}

type UserRole struct {
	UserID     uint      `gorm:"primaryKey;autoIncrement:false" json:"user_id"`
	RoleID     uint      `gorm:"primaryKey;autoIncrement:false" json:"role_id"`
	TenantID   uint      `gorm:"not null;index" json:"tenant_id"`
	AssignedBy uint      `json:"assigned_by,omitzero"`
	AssignedAt time.Time `gorm:"autoCreateTime" json:"assigned_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
	Role Role `gorm:"foreignKey:RoleID" json:"-"`
}

type Department struct {
	ID        uint           `gorm:"primaryKey"`
	TenantID  uint           `gorm:"not null;index"`
	Name      string         `gorm:"size:100;not null"`
	Code      string         `gorm:"size:20;uniqueIndex:idx_dept_tenant_code,composite:tenant_code"`
	IsActive  bool           `gorm:"default:true"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Positions []Position     `gorm:"foreignKey:DepartmentID"`
}

type Position struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     uint           `gorm:"not null;index" json:"tenant_id"`
	DepartmentID uint           `gorm:"index"`
	NamePosition string         `gorm:"size:100;not null" json:"name_position"`
	Description  string         `gorm:"size:255;not null" json:"description"`
	IsActive     bool           `gorm:"default:true"`
	CreatedAt    time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index:idx_users_deleted_at" json:"deleted_at,omitzero"`
	Deparment    Department     `gorm:"foreignKey:DepartmentID"`
}

type Employee struct {
	ID           uint   `gorm:"primaryKey"`
	TenantID     uint   `gorm:"not null;index"`
	UserID       uint   `gorm:"not null;uniqueIndex"`
	DepartmentID uint   `gorm:"index"`
	PositionID   uint   `gorm:"index"`
	IsActive     bool   `gorm:"default:true;index"`
	EmployeeCode string `gorm:"size:20;index"`            // Código interno del empleado
	Status       string `gorm:"size:20;default:'active'"` // active, inactive, suspended
	Notes        string `gorm:"size:500"`                 // Notas internas
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	User       User               `gorm:"foreignKey:UserID"`
	Department Department         `gorm:"foreignKey:DepartmentID"`
	Position   Position           `gorm:"foreignKey:PositionID"`
	Contracts  []EmployeeContract `gorm:"foreignKey:EmployeeID"`
}
type EmployeeContract struct {
	ID             uint `gorm:"primaryKey"`
	TenantID       uint `gorm:"not null;index"`
	EmployeeID     uint `gorm:"not null;index"`
	ContractTypeID uint `gorm:"index"`

	BaseSalary float64
	Currency   string `gorm:"size:3"`

	StartDate time.Time
	EndDate   *time.Time
	IsActive  bool `gorm:"index"`

	WorkHoursPerDay     float64
	WorkDaysPerWeek     float64
	HealthContribution  float64
	PensionContribution float64
	TransportAllowance  float64
	HousingAllowance    float64

	CreatedAt time.Time
	UpdatedAt time.Time

	Employee     Employee     `gorm:"foreignKey:EmployeeID"`
	ContractType ContractType `gorm:"foreignKey:ContractTypeID"`
}

type ContractType struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:50;not null"`
	Description string `gorm:"size:255"`
}

type Payroll struct {
	ID         uint `gorm:"primaryKey"`
	TenantID   uint `gorm:"not null;index"`
	EmployeeID uint `gorm:"not null;index"`

	PeriodStart     time.Time
	PeriodEnd       time.Time
	PayDate         time.Time
	GrossAmount     float64
	TotalDeductions float64
	NetAmount       float64
	Status          string `gorm:"size:20;default:'draft'"` // draft, calculated, paid
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Employee Employee      `gorm:"foreignKey:EmployeeID"`
	Items    []PayrollItem `gorm:"foreignKey:PayrollID"`
}

type PayrollItem struct {
	ID           uint    `gorm:"primaryKey"`
	PayrollID    uint    `gorm:"not null;index"`
	ConceptID    uint    `gorm:"index"`
	Type         string  `gorm:"size:20"`       // earning | deduction | employer_contribution
	Code         string  `gorm:"size:30;index"` // SALARY, HEALTH_EMPLOYEE, PENSION_EMPLOYER, TAX
	Name         string  `gorm:"size:100"`
	Amount       float64 `gorm:"not null"`
	CalculatedAt time.Time

	Payroll Payroll        `gorm:"foreignKey:PayrollID"`
	Concept PayrollConcept `gorm:"foreignKey:ConceptID"`
}

type PayrollConcept struct {
	ID           uint           `gorm:"primaryKey"`
	TenantID     uint           `gorm:"not null;index" json:"tenant_id"`
	Code         string         `gorm:"size:30;not null;uniqueIndex:idx_concept_tenant_code,composite:tenant_code" json:"code"`
	Name         string         `gorm:"size:100" json:"name"`
	Type         string         `gorm:"size:20" json:"type"` // earning | deduction | employer_contribution
	Description  string         `gorm:"size:255" json:"description"`
	Percentage   float64        `gorm:"default:0" json:"percentage"`
	EmployeePart float64        `gorm:"default:0" json:"employee_part"`
	EmployerPart float64        `gorm:"default:0" json:"employer_part"`
	IsMandatory  bool           `gorm:"default:false" json:"is_mandatory"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`
}

type Payment struct {
	ID        uint `gorm:"primaryKey"`
	PayrollID uint `gorm:"not null;index"`

	Method string `gorm:"size:30"` // bank_transfer

	BankName      string `gorm:"size:100"`
	AccountNumber string `gorm:"size:50"`

	Amount float64

	PaidAt    time.Time
	Status    string `gorm:"size:20"`
	CreatedAt time.Time

	Payroll Payroll `gorm:"foreignKey:PayrollID"`
}

// ========================================
// Horas Extras
// ========================================

const (
	OvertimeTypeExtra      = "extra"       // Hora extra diurna
	OvertimeTypeNight      = "night"       // Recargo nocturno
	OvertimeTypeHoliday    = "holiday"     // Hora en día festivo
	OvertimeTypeSunday     = "sunday"      // Recargo dominical
	OvertimeTypeExtraNight = "extra_night" // Hora extra nocturna
)

// Overtime representa horas extras de un empleado
type Overtime struct {
	ID         uint      `gorm:"primaryKey"`
	TenantID   uint      `gorm:"not null;index"`
	EmployeeID uint      `gorm:"not null;index"`
	PayrollID  uint      `gorm:"index"` // nullable, se asigna al procesar nómina
	Date       time.Time `gorm:"not null;index"`
	Hours      float64   `gorm:"not null"`
	Type       string    `gorm:"size:20;not null"` // extra, night, holiday, sunday
	Rate       float64   `gorm:"default:1.0"`      // multiplicador (1.25, 1.35, 2.0)
	Amount     float64   `gorm:"default:0"`        // calculado
	ApprovedBy uint      `gorm:"index"`
	ApprovedAt *time.Time
	Status     string `gorm:"size:20;default:'pending'"` // pending, approved, rejected
	Notes      string `gorm:"size:500"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	Employee Employee `gorm:"foreignKey:EmployeeID"`
	Payroll  Payroll  `gorm:"foreignKey:PayrollID"`
}

func (Overtime) TableName() string {
	return "overtimes"
}

// ========================================
// Incapacidades / Licencias
// ========================================

const (
	AbsenceTypeSick      = "sick"      // Incapacidad médica
	AbsenceTypeMaternity = "maternity" // Licencia de maternidad
	AbsenceTypePaternity = "paternity" // Licencia de paternidad
	AbsenceTypeVacation  = "vacation"  // Vacaciones
	AbsenceTypeUnpaid    = "unpaid"    // Licencia sin pago
	AbsenceTypeOther     = "other"     // Otra ausencia
)

// Absence representa una ausencia/incapacidad
type Absence struct {
	ID                 uint      `gorm:"primaryKey"`
	TenantID           uint      `gorm:"not null;index"`
	EmployeeID         uint      `gorm:"not null;index"`
	PayrollID          uint      `gorm:"index"`            // nullable
	Type               string    `gorm:"size:20;not null"` // sick, vacation, maternity, etc.
	StartDate          time.Time `gorm:"not null;index"`
	EndDate            time.Time `gorm:"not null"`
	PaidPercent        float64   `gorm:"default:100"`              // porcentaje de pago (100, 66.6, 0)
	DailyRate          float64   `gorm:"default:0"`                // valor diario
	TotalDeduction     float64   `gorm:"default:0"`                // total deduccion
	MedicalCertificate string    `gorm:"size:50"`                  // certificado médico
	Status             string    `gorm:"size:20;default:'active'"` // active, processed
	Notes              string    `gorm:"size:500"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`

	Employee Employee `gorm:"foreignKey:EmployeeID"`
	Payroll  Payroll  `gorm:"foreignKey:PayrollID"`
}

func (Absence) TableName() string {
	return "absences"
}

// ========================================
// Bonificaciones
// ========================================

const (
	BonusTypePerformance = "performance" // Por desempeño
	BonusTypeProduction  = "production"  // Por producción
	BonusTypeAttendance  = "attendance"  // Por asistencia
	BonusTypeChristmas   = "christmas"   // Prima de servicios / navidad
	BonusTypeOther       = "other"       // Otra bonificación
)

// Bonus representa una bonificación variable
type Bonus struct {
	ID          uint      `gorm:"primaryKey"`
	TenantID    uint      `gorm:"not null;index"`
	EmployeeID  uint      `gorm:"not null;index"`
	PayrollID   uint      `gorm:"index"`
	Type        string    `gorm:"size:20;not null"` // performance, production, etc.
	Amount      float64   `gorm:"not null"`
	Description string    `gorm:"size:255"`
	Date        time.Time `gorm:"not null;index"` // fecha de la bonificación
	ApprovedBy  uint      `gorm:"index"`
	ApprovedAt  *time.Time
	Status      string `gorm:"size:20;default:'pending'"` // pending, approved, paid
	Notes       string `gorm:"size:500"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	Employee Employee `gorm:"foreignKey:EmployeeID"`
	Payroll  Payroll  `gorm:"foreignKey:PayrollID"`
}

func (Bonus) TableName() string {
	return "bonuses"
}

// ========================================
// Retención en Fuente (Impuestos)
// ========================================

// TaxRule representa una regla de retención
type TaxRule struct {
	ID          uint    `gorm:"primaryKey"`
	TenantID    uint    `gorm:"not null;index"`
	Code        string  `gorm:"size:30;not null;uniqueIndex:idx_tax_rule_tenant_code,composite:tenant_code"`
	Name        string  `gorm:"size:100;not null"`
	BasePercent float64 `gorm:"default:0"` // porcentaje base de retención
	TaxCategory string  `gorm:"size:30"`   // ingresos_trabajo, honorarios, etc.
	MinIncome   float64 `gorm:"default:0"` // ingreso mínimo para aplicar
	MaxIncome   float64 `gorm:"default:0"` // ingreso máximo (0 = sin límite)
	Priority    int     `gorm:"default:0"` // orden de aplicación
	IsActive    bool    `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (TaxRule) TableName() string {
	return "tax_rules"
}

// ========================================
// Métodos de conveniencia
// ========================================

// TableName especifica nombres de tablas
func (User) TableName() string {
	return "users"
}

func (Permission) TableName() string {
	return "permissions"
}

func (PermissionAction) TableName() string {
	return "permission_actions"
}

func (Role) TableName() string {
	return "roles"
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func (UserRole) TableName() string {
	return "user_roles"
}

func (Department) TableName() string {
	return "departments"
}

func (Position) TableName() string {
	return "positions"
}

func (Employee) TableName() string {
	return "employees"
}

func (EmployeeContract) TableName() string {
	return "employee_contracts"
}

func (ContractType) TableName() string {
	return "contract_types"
}

func (Payroll) TableName() string {
	return "payrolls"
}

func (PayrollItem) TableName() string {
	return "payroll_items"
}

func (PayrollConcept) TableName() string {
	return "payroll_concepts"
}

func (Payment) TableName() string {
	return "payments"
}

// ========================================
// Métodos de User para verificar permisos
// ========================================

// HasRole verifica si el usuario tiene un rol específico
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName && role.IsActive {
			return true
		}
	}
	return false
}

// HasPermission verifica si el usuario tiene un permiso específico
func (u *User) HasPermission(resource, action string) bool {
	for _, role := range u.Roles {
		if !role.IsActive {
			continue
		}
		for _, rp := range role.RolePermissions {
			if rp.Action.Resource.Name == resource &&
				rp.Action.Action == action &&
				rp.Action.IsActive {
				return true
			}
		}
	}
	return false
}

// IsAdmin verifica si el usuario es admin
func (u *User) IsAdmin() bool {
	return u.HasRole("admin")
}

// GetAllPermissions retorna todos los permisos del usuario
func (u *User) GetAllPermissions() []string {
	permissions := make(map[string]bool)

	for _, role := range u.Roles {
		if !role.IsActive {
			continue
		}
		for _, rp := range role.RolePermissions {
			if rp.Action.IsActive {
				slug := rp.Action.Resource.Name + "." + rp.Action.Action
				permissions[slug] = true
			}
		}
	}

	result := make([]string, 0, len(permissions))
	for perm := range permissions {
		result = append(result, perm)
	}
	return result
}

// ========================================
// Métodos de Role
// ========================================

// HasPermission verifica si el rol tiene un permiso específico
func (r *Role) HasPermission(resource, action string) bool {
	for _, rp := range r.RolePermissions {
		if rp.Action.Resource.Name == resource &&
			rp.Action.Action == action &&
			rp.Action.IsActive {
			return true
		}
	}
	return false
}

// ========================================
// Métodos de Employee (Behaviors)
// ========================================

// HasActiveContract verifica si el empleado tiene contrato activo
func (e *Employee) HasActiveContract() bool {
	for _, c := range e.Contracts {
		if c.IsActive {
			return true
		}
	}
	return false
}

// GetMainContract retorna el contrato activo principal
func (e *Employee) GetMainContract() *EmployeeContract {
	for _, c := range e.Contracts {
		if c.IsActive {
			return &c
		}
	}
	return nil
}

// GetFullName retorna el nombre completo del empleado
func (e *Employee) GetFullName() string {
	if e.User.FirstName == "" && e.User.LastName == "" {
		return ""
	}
	return e.User.FirstName + " " + e.User.LastName
}

// CanBeActivated verifica si el empleado puede ser activado
func (e *Employee) CanBeActivated() bool {
	return e.Status == EmployeeStatusInactive || e.Status == EmployeeStatusSuspended
}

// CanBeDeactivated verifica si el empleado puede ser desactivado
func (e *Employee) CanBeDeactivated() bool {
	// No puede desactivarse si tiene contrato activo
	for _, c := range e.Contracts {
		if c.IsActive {
			return false
		}
	}
	return true
}

// SetStatus cambia el estado del empleado con validaciones
func (e *Employee) SetStatus(newStatus string) error {
	oldStatus := e.Status

	switch newStatus {
	case EmployeeStatusActive:
		if oldStatus != EmployeeStatusInactive && oldStatus != EmployeeStatusSuspended {
			return errors.New("can only activate inactive or suspended employees")
		}
	case EmployeeStatusInactive:
		if oldStatus == EmployeeStatusActive && !e.CanBeDeactivated() {
			return errors.New("cannot deactivate: employee has active contract")
		}
	case EmployeeStatusSuspended:
		// Suspension es reversible
		e.Status = newStatus
		return nil
	}

	e.Status = newStatus
	return nil
}

// Validate valida los datos del empleado
func (e *Employee) Validate() error {
	if e.UserID == 0 {
		return errors.New("user is required")
	}
	if e.TenantID == 0 {
		return errors.New("tenant is required")
	}
	if e.Status != "" && e.Status != EmployeeStatusActive &&
		e.Status != EmployeeStatusInactive && e.Status != EmployeeStatusSuspended {
		return errors.New("invalid status: must be active, inactive, or suspended")
	}
	return nil
}

// ========================================
// Métodos de PermissionAction
// ========================================

// GetSlug retorna el slug del permiso (users.create)
func (pa *PermissionAction) GetSlug() string {
	return pa.Resource.Name + "." + pa.Action
}
