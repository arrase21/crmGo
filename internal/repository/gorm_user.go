package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/arrase21/crm-users/internal/domain"
	"gorm.io/gorm"
)

type GormUserRepo struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) domain.UserRepo {
	return &GormUserRepo{
		db: db,
	}
}

func tenantFromctx(ctx context.Context) (uint, error) {
	tenantID, ok := ctx.Value(domain.TenantIDKey).(uint)
	if !ok || tenantID == 0 {
		return 0, errors.New("tenant not found in context")
	}
	return tenantID, nil
}

func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "duplicate") ||
		strings.Contains(errMsg, "unique constraint")
}

func (r *GormUserRepo) Create(ctx context.Context, usr *domain.User) error {
	if usr == nil {
		return errors.New("user cannot be nil")
	}

	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}

	usr.TenantID = tenantID

	err = r.db.WithContext(ctx).Create(usr).Error
	if err != nil {
		if isDuplicateError(err) {
			if strings.Contains(err.Error(), "dni") {
				return domain.ErrDniAlreadyExist
			}
			if strings.Contains(err.Error(), "email") {
				return domain.ErrEmailAlreadyExist
			}
			if strings.Contains(err.Error(), "phone") {
				return domain.ErrPhoneAlreadyExist
			}
			return errors.New("duplicate entry")
		}
		return err
	}

	return nil
}

func (r *GormUserRepo) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user id")
	}

	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *GormUserRepo) GetByDni(ctx context.Context, dni string) (*domain.User, error) {
	if dni == "" {
		return nil, errors.New("dni cannot be empty")
	}

	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = r.db.WithContext(ctx).
		Where("tenant_id = ? AND dni = ?", tenantID, dni).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *GormUserRepo) List(ctx context.Context) ([]domain.User, error) {
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}

	var usrs []domain.User
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("id DESC").
		Find(&usrs).Error; err != nil {
		return nil, err
	}

	return usrs, nil
}

func (r *GormUserRepo) Update(ctx context.Context, usr *domain.User) error {
	if usr == nil || usr.ID == 0 {
		return errors.New("user cannot be nil or have zero id")
	}

	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}

	// Verificar que el usuario existe y pertenece al tenant
	existing, err := r.GetByID(ctx, usr.ID)
	if err != nil {
		return err
	}

	// Preservar tenant_id original (seguridad)
	usr.TenantID = existing.TenantID

	err = r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ? AND tenant_id = ?", usr.ID, tenantID).
		Updates(usr).Error

	if err != nil {
		if isDuplicateError(err) {
			if strings.Contains(err.Error(), "dni") {
				return domain.ErrDniAlreadyExist
			}
			if strings.Contains(err.Error(), "email") {
				return domain.ErrEmailAlreadyExist
			}
			if strings.Contains(err.Error(), "phone") {
				return domain.ErrPhoneAlreadyExist
			}
			return errors.New("duplicate entry")
		}
		return err
	}

	return nil
}

func (r *GormUserRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid user id")
	}

	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&domain.User{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
