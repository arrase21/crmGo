package repository

import (
	"context"
	"errors"

	"github.com/arrase21/crm-users/internal/domain"
	"gorm.io/gorm"
)

type GormPayrollItemRepo struct {
	db *gorm.DB
}

func NewGormPayrollItemRepository(db *gorm.DB) domain.PayrollItemRepo {
	return &GormPayrollItemRepo{
		db: db,
	}
}

func (r *GormPayrollItemRepo) Create(ctx context.Context, item *domain.PayrollItem) error {
	if item == nil {
		return errors.New("item cannot be nil")
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *GormPayrollItemRepo) CreateBatch(ctx context.Context, items []domain.PayrollItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *GormPayrollItemRepo) GetByIDPayrollID(ctx context.Context, payrollID uint) ([]domain.PayrollItem, error) {
	if payrollID == 0 {
		return nil, errors.New("invalid payrollid")
	}
	var items []domain.PayrollItem
	err := r.db.WithContext(ctx).Where("payroll_id = ?", payrollID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, err
}

func (r *GormPayrollItemRepo) DeleteByPayrollID(ctx context.Context, payrollID uint) error {
	if payrollID == 0 {
		return errors.New("payroll cannot be nil")
	}
	result := r.db.WithContext(ctx).Where("payroll_id = ?", payrollID).Delete(&domain.PayrollItem{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no items to delete")
	}
	return nil
}
