package domain

import (
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
	ID           uint `gorm:"primaryKey"`
	TenantID     uint `gorm:"not null;index"`
	UserID       uint `gorm:"not null;uniqueIndex"`
	DepartmentID uint `gorm:"index"`
	PositionID   uint `gorm:"index"`
	IsActive     bool `gorm:"default:true;index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt     `gorm:"index"`
	User         User               `gorm:"foreignKey:UserID"`
	Department   Department         `gorm:"foreignKey:DepartmentID"`
	Position     Position           `gorm:"foreignKey:PositionID"`
	Contracts    []EmployeeContract `gorm:"foreignKey:EmployeeID"`
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
// Métodos de PermissionAction
// ========================================

// GetSlug retorna el slug del permiso (users.create)
func (pa *PermissionAction) GetSlug() string {
	return pa.Resource.Name + "." + pa.Action
}
