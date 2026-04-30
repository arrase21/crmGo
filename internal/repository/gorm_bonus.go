package repository

import (
	"context"
	"errors"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
	"gorm.io/gorm"
)

type GormBonusRepo struct {
	db *gorm.DB
}

func NewGormBonusRepository(db *gorm.DB) domain.BonusRepo {
	return &GormBonusRepo{db: db}
}

func (r *GormBonusRepo) Create(ctx context.Context, bonus *domain.Bonus) error {
	if bonus == nil {
		return errors.New("bonus cannot be nil")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	bonus.TenantID = tenantID
	return r.db.WithContext(ctx).Create(bonus).Error
}

func (r *GormBonusRepo) GetByID(ctx context.Context, id uint) (*domain.Bonus, error) {
	if id == 0 {
		return nil, errors.New("invalid bonus id")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var bonus domain.Bonus
	err = r.db.WithContext(ctx).
		Preload("Employee.User").
		Where("tenant_id = ? AND id = ?", tenantID, id).First(&bonus).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bonus not found")
		}
		return nil, err
	}
	return &bonus, nil
}

func (r *GormBonusRepo) GetByEmployeeAndPeriod(ctx context.Context, employeeID uint, start, end time.Time) ([]domain.Bonus, error) {
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var bonuses []domain.Bonus
	err = r.db.WithContext(ctx).
		Where("tenant_id = ? AND employee_id = ? AND date >= ? AND date <= ?",
			tenantID, employeeID, start, end).
		Order("date ASC").
		Find(&bonuses).Error
	return bonuses, err
}

func (r *GormBonusRepo) GetByPayrollID(ctx context.Context, payrollID uint) ([]domain.Bonus, error) {
	var bonuses []domain.Bonus
	err := r.db.WithContext(ctx).
		Where("payroll_id = ?", payrollID).
		Find(&bonuses).Error
	return bonuses, err
}

func (r *GormBonusRepo) List(ctx context.Context, employeeID uint, page, limit int) ([]domain.Bonus, int64, error) {
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var bonuses []domain.Bonus
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Bonus{}).Where("tenant_id = ?", tenantID)
	if employeeID != 0 {
		query = query.Where("employee_id = ?", employeeID)
	}

	query.Count(&total)

	err = query.Preload("Employee.User").
		Order("date DESC").
		Offset(offset).
		Limit(limit).
		Find(&bonuses).Error

	return bonuses, total, err
}

func (r *GormBonusRepo) Update(ctx context.Context, bonus *domain.Bonus) error {
	if bonus == nil || bonus.ID == 0 {
		return errors.New("bonus cannot be nil or without id")
	}
	return r.db.WithContext(ctx).Save(bonus).Error
}

func (r *GormBonusRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid bonus id")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&domain.Bonus{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("bonus not found")
	}
	return nil
}
