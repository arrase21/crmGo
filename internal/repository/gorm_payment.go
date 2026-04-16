package repository

import (
	"context"
	"errors"

	"github.com/arrase21/crm-users/internal/domain"
	"gorm.io/gorm"
)

type GormPaymentRepo struct {
	db *gorm.DB
}

func NewGormPaymentRepository(db *gorm.DB) domain.PaymentRepo {
	return &GormPaymentRepo{db: db}
}

func (r *GormPaymentRepo) Create(ctx context.Context, payment *domain.Payment) error {
	if payment == nil {
		return errors.New("payment cannot be nil")
	}
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *GormPaymentRepo) GetByID(ctx context.Context, id uint) (*domain.Payment, error) {
	if id == 0 {
		return nil, errors.New("invalid payment id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}

	var payment domain.Payment
	err = r.db.WithContext(ctx).
		Preload("Payroll").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}
	return &payment, nil
}

func (r *GormPaymentRepo) GetByPayrollID(ctx context.Context, payrollID uint) (*domain.Payment, error) {
	if payrollID == 0 {
		return nil, errors.New("invalid payroll id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}

	var payment domain.Payment
	err = r.db.WithContext(ctx).
		Where("payroll_id = ?", payrollID).
		First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}
	_ = tenantID // Evitar warning de variable no usada
	return &payment, nil
}

func (r *GormPaymentRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid payment id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&domain.Payment{}).Error
}
