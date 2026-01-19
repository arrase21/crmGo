package domain

import (
	"time"

	"gorm.io/gorm"
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
