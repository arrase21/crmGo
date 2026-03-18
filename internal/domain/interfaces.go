package domain

import (
	"context"
	"errors"
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
