package repository

import (
	"context"
	"errors"

	"github.com/arrase21/crm-users/internal/domain"
	"gorm.io/gorm"
)

type GormUserRepo struct {
	db       *gorm.DB
	tenantID uint
}

func NewGormUserRepository(db *gorm.DB) domain.UserRepo {
	return &GormUserRepo{
		db: db,
	}
}

func tenantFromctx(ctx context.Context) (uint, error) {
	tenantID, ok := ctx.Value(domain.TenantIDKey).(uint)
	if !ok || tenantID == 0 {
		return 0, errors.New("tenant not found")
	}
	return tenantID, nil
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
	return r.db.WithContext(ctx).Create(usr).Error
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
	err = r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Where("id = ?", id).First(&user).Error
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
	err = r.db.WithContext(ctx).Where("tenant_id = ? AND dni = ?", tenantID, dni).First(&user).Error
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
		Find(&usrs).Error; err != nil {
		return nil, err
	}
	return usrs, nil
}

func (r *GormUserRepo) Update(ctx context.Context, usr *domain.User) error {
	if usr == nil || usr.ID == 0 {
		return errors.New("user cant be nil or 0")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}

	existing, err := r.GetByID(ctx, usr.ID)
	if err != nil {
		return err
	}
	usr.TenantID = existing.TenantID

	result := r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ? AND tenant_id = ?", usr.ID, tenantID).
		Updates(usr)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *GormUserRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid user id")
	}

	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil
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
