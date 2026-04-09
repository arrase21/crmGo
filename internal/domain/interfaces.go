package domain

import (
	"context"
	"errors"
	"time"
)

// domain errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDniAlreadyExist   = errors.New("dni already exists")
	ErrEmailAlreadyExist = errors.New("email already exists")
	ErrPhoneAlreadyExist = errors.New("phone already exists")
	ErrTenantNotFound    = errors.New("tenant not found in context")
	ErrInvalidTenantID   = errors.New("invalid tenant id")
)

// Errores de dominio - Roles
var (
	ErrRoleNotFound = errors.New("role not found")
	ErrRoleExisting = errors.New("role already exists")
)

// Errores de dominio - Permissions
var (
	ErrPermissionNotFound = errors.New("permission not found")
	ErrActionNotFound     = errors.New("action not found")
)

// Errores de Nómina
var (
	ErrEmployeeNotFound         = errors.New("employee not found")
	ErrEmployeeContractNotFound = errors.New("active contract not found")
	ErrPayrollNotFound          = errors.New("payroll not found")
	ErrPayrollAlreadyPaid       = errors.New("payroll already paid")
	ErrConceptNotFound          = errors.New("payroll concept not found")
)

// ContextKey for tenant
type contextKey string

const TenantIDKey contextKey = "tenant_id"

type UserRepo interface {
	Create(ctx context.Context, usr *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByDni(ctx context.Context, dni string) (*User, error)
	List(ctx context.Context, page, limit int) ([]User, int64, error)
	Update(ctx context.Context, usr *User) error
	Delete(ctx context.Context, id uint) error
}

type RoleRepo interface {
	Create(ctx context.Context, role *Role) error
	GetByID(ctx context.Context, roleID uint) (*Role, error)
	GetByName(ctx context.Context, name string) (*Role, error)
	List(ctx context.Context) ([]Role, error)
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, roleID uint) error

	// Gestión de permisos
	AssignPermission(ctx context.Context, roleID, actionID uint) error
	RevokePermission(ctx context.Context, roleID, actionID uint) error
	GetPermissions(ctx context.Context, roleID uint) ([]PermissionAction, error)
}
type PermissionRepo interface {
	CreatePermission(ctx context.Context, perm *Permission) error
	GetPermissionByID(ctx context.Context, id uint) (*Permission, error)
	GetPermissionByName(ctx context.Context, name string) (*Permission, error)
	ListPermission(ctx context.Context) ([]Permission, error)
	//Actions
	CreateAction(ctx context.Context, action *PermissionAction) error
	GetActionByID(ctx context.Context, id uint) (*PermissionAction, error)
	ListActions(ctx context.Context, resourceID uint) ([]PermissionAction, error)
	ListAllActions(ctx context.Context) ([]PermissionAction, error)
}

type UserRoleRepo interface {
	AssignRole(ctx context.Context, userID, roleID uint) error
	RevokeRole(ctx context.Context, userID, roleID uint) error
	GetUserRoles(ctx context.Context, userID uint) ([]Role, error)
	GetRoleUsers(ctx context.Context, userID uint) ([]User, error)
}

// Los siguientes modelos estan para implementar al acabar de implementar uno bajar este al siguiente modelo
type EmployeeRepo interface {
	Create(ctx context.Context, emp *Employee) error
	GetByID(ctx context.Context, id uint) (*Employee, error)
	GetByUserID(ctx context.Context, userID uint) (*Employee, error)
	List(ctx context.Context, page, limit int) ([]Employee, int64, error)
	Update(ctx context.Context, emp *Employee) error
	Delete(ctx context.Context, id uint) error
}

type PayrollRepo interface {
	Create(ctx context.Context, payroll *Payroll) error
	GetByID(ctx context.Context, id uint) (*Payroll, error)
	GetByEmployeeAndPeriod(ctx context.Context, employeeID uint, periodStart, periodEnd time.Time) (*Payroll, error)
	ListByEmployee(ctx context.Context, employeeID uint) ([]Payroll, error)
	Update(ctx context.Context, payroll *Payroll) error
	Delete(ctx context.Context, id uint) error
}

type PayrollConceptRepo interface {
	Create(ctx context.Context, concept *PayrollConcept) error
	GetByID(ctx context.Context, id uint) (*PayrollConcept, error)
	GetByCode(ctx context.Context, code string) (*PayrollConcept, error)
	GetActiveConcepts(ctx context.Context) ([]PayrollConcept, error)
	List(ctx context.Context, page, limit int) ([]PayrollConcept, int64, error)
	Update(ctx context.Context, concept *PayrollConcept) error
	Delete(ctx context.Context, id uint) error
}
