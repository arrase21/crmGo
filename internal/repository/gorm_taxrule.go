package repository

import (
	"context"
	"errors"

	"github.com/arrase21/crm-users/internal/domain"
	"gorm.io/gorm"
)

type GormTaxRuleRepo struct {
	db *gorm.DB
}

func NewGormTaxRuleRepository(db *gorm.DB) domain.TaxRuleRepo {
	return &GormTaxRuleRepo{db: db}
}

func (r *GormTaxRuleRepo) Create(ctx context.Context, rule *domain.TaxRule) error {
	if rule == nil {
		return errors.New("tax rule cannot be nil")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	rule.TenantID = tenantID
	return r.db.WithContext(ctx).Create(rule).Error
}

func (r *GormTaxRuleRepo) GetByID(ctx context.Context, id uint) (*domain.TaxRule, error) {
	if id == 0 {
		return nil, errors.New("invalid tax rule id")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var rule domain.TaxRule
	err = r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).First(&rule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tax rule not found")
		}
		return nil, err
	}
	return &rule, nil
}

func (r *GormTaxRuleRepo) GetByIncome(ctx context.Context, income float64) (*domain.TaxRule, error) {
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var rule domain.TaxRule
	err = r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_active = ? AND min_income <= ?", tenantID, true, income).
		Order("priority DESC").
		First(&rule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no tax rule found for this income")
		}
		return nil, err
	}
	return &rule, nil
}

func (r *GormTaxRuleRepo) List(ctx context.Context, page, limit int) ([]domain.TaxRule, int64, error) {
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

	var rules []domain.TaxRule
	var total int64

	err = r.db.WithContext(ctx).Model(&domain.TaxRule{}).
		Where("tenant_id = ?", tenantID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("priority ASC").
		Offset(offset).
		Limit(limit).
		Find(&rules).Error

	return rules, total, err
}

func (r *GormTaxRuleRepo) Update(ctx context.Context, rule *domain.TaxRule) error {
	if rule == nil || rule.ID == 0 {
		return errors.New("tax rule cannot be nil or without id")
	}
	return r.db.WithContext(ctx).Save(rule).Error
}

func (r *GormTaxRuleRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid tax rule id")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&domain.TaxRule{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("tax rule not found")
	}
	return nil
}
