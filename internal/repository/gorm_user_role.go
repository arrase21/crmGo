package repository

import (
	"context"
	"errors"

	"github.com/arrase21/crm-users/internal/domain"
	"gorm.io/gorm"
)

type GormUserRoleRepo struct {
	db *gorm.DB
}

func NewGormUserRoleRepository(db *gorm.DB) domain.UserRoleRepo {
	return &GormUserRoleRepo{db: db}
}

func (r *GormUserRoleRepo) AssignRole(ctx context.Context, userID, roleID uint) error {
	if userID == 0 || roleID == 0 {
		return errors.New("cannot be 0")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	var user domain.User
	if err := r.db.Where("id = ? AND tenant_id = ?", userID, tenantID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		return err
	}
	if user.TenantID != tenantID {
		return errors.New("user does not belong to tenant")
	}
	var role domain.Role
	if err := r.db.Where("id = ? AND tenant_id = ?", roleID, tenantID).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrRoleNotFound
		}
		return err
	}
	if role.TenantID != tenantID {
		return errors.New("role does not belong to tenant")
	}
	userRole := &domain.UserRole{
		UserID:   userID,
		RoleID:   roleID,
		TenantID: tenantID,
	}
	err = r.db.WithContext(ctx).Create(userRole).Error
	if err != nil {
		if isDuplicateError(err) {
			return nil
		}
		return err
	}
	return nil
}

func (r *GormUserRoleRepo) RevokeRole(ctx context.Context, userID, roleID uint) error {
	if userID == 0 || roleID == 0 {
		return errors.New("invalid userid or roleid")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).Where("user_id = ? AND role_id = ? AND tenant_id = ?", userID, roleID, tenantID).
		Delete(&domain.UserRole{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user role assignment not found")
	}
	return nil
}

func (r *GormUserRoleRepo) GetUserRoles(ctx context.Context, userID uint) ([]domain.Role, error) {
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var roles []domain.Role
	err = r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Preload("RolePermissions.Action.Resource").
		Where("user_roles.user_id = ? AND user_roles.tenant_id = ?", userID, tenantID).
		Where("roles.is_active = ?", true).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *GormUserRoleRepo) GetRoleUsers(ctx context.Context, roleID uint) ([]domain.User, error) {
	if roleID == 0 {
		return nil, errors.New("invalid user id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var users []domain.User
	err = r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ? AND user_roles.tenant_id = ?", roleID, tenantID).
		Find(&users).Error
	if err != nil {
		return nil, errors.New("invalid role id")
	}
	return users, nil
}
