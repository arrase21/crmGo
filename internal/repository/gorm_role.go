package repository

import (
	"context"
	"errors"

	"github.com/arrase21/crm-users/internal/domain"
	"gorm.io/gorm"
)

type GormRoleRepo struct {
	db *gorm.DB
}

func NewGormRoleRepository(db *gorm.DB) domain.RoleRepo {
	return &GormRoleRepo{
		db: db,
	}
}

func (r *GormRoleRepo) Create(ctx context.Context, role *domain.Role) error {
	if role == nil {
		return errors.New("role cannot be nil")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}

	if role.TenantID == 0 {
		role.TenantID = tenantID
	}

	err = r.db.WithContext(ctx).Create(role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormRoleRepo) GetByID(ctx context.Context, id uint) (*domain.Role, error) {
	if id == 0 {
		return nil, errors.New("invalid role id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var role domain.Role
	err = r.db.WithContext(ctx).
		Preload("RolePermissions.Action.Resource").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

func (r *GormRoleRepo) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var role domain.Role
	err = r.db.WithContext(ctx).
		Preload("RolePermissions.Action.Resource").
		Where("tenant_id = ? AND name = ?", tenantID, name).
		First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

func (r *GormRoleRepo) List(ctx context.Context) ([]domain.Role, error) {
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var roles []domain.Role
	if err := r.db.WithContext(ctx).
		Preload("RolePermissions.Action.Resource").
		Where("tenant_id = ?", tenantID).
		Order("id DESC").
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *GormRoleRepo) Update(ctx context.Context, role *domain.Role) error {
	if role == nil || role.ID == 0 {
		return errors.New("role cannot be nil or have zero id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	existing, err := r.GetByID(ctx, role.ID)
	if err != nil {
		return err
	}

	role.TenantID = existing.TenantID

	err = r.db.WithContext(ctx).
		Model(&domain.Role{}).
		Where("id = ? AND tenant_id = ?", role.ID, tenantID).Updates(role).Error
	if err != nil {
		if isDuplicateError(err) {
			return domain.ErrRoleExisting
		}
		return err
	}
	return nil
}

func (r *GormRoleRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}

	role, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if role.IsSystem {
		return errors.New("cannot delete system roles")
	}

	result := r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.Role{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrRoleNotFound
	}
	return nil
}

// Assing permissions
func (r *GormRoleRepo) AssignPermission(ctx context.Context, roleID, actionID uint) error {
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	role, err := r.GetByID(ctx, roleID)
	if err != nil {
		return err
	}

	if role.TenantID != tenantID {
		return errors.New("role does not belong to tenant")
	}

	var action domain.PermissionAction
	if err := r.db.First(&action, actionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrActionNotFound
		}
		return err
	}
	rolePermission := &domain.RolePermission{
		RoleID:   roleID,
		ActionID: actionID,
	}
	err = r.db.WithContext(ctx).Create(rolePermission).Error
	if err != nil {
		if isDuplicateError(err) {
			return nil
		}
		return err
	}
	return nil
}

func (r *GormRoleRepo) RevokePermission(ctx context.Context, roleID, actionID uint) error {
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	role, err := r.GetByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role.TenantID != tenantID {
		return errors.New("role does not belong to tenant")
	}
	result := r.db.WithContext(ctx).Where("role_is = ? AND action_id = ?", roleID, actionID).Delete(&domain.RolePermission{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormRoleRepo) GetPermissions(ctx context.Context, roleID uint) ([]domain.PermissionAction, error) {
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}
	role, err := r.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role.TenantID != tenantID {
		return nil, errors.New("role does not belong to tenant")
	}
	var actions []domain.PermissionAction
	err = r.db.WithContext(ctx).Joins("JOIN role_permissions ON role_permissions.action_id = permission_actions.id").
		Preload("Resource").
		Where("role_permissions.role_id = ?", roleID).
		Find(&actions).Error

	if err != nil {
		return nil, err
	}
	return actions, nil
}
