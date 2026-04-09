package repository

import (
	"context"
	"errors"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
	"gorm.io/gorm"
)

type GormPayrollRepo struct {
	db *gorm.DB
}

func NewGormPayrollRepository(db *gorm.DB) domain.PayrollRepo {
	return &GormPayrollRepo{
		db: db,
	}
}

func (r *GormPayrollRepo) Create(ctx context.Context, payroll *domain.Payroll) error {
	if payroll == nil {
		return errors.New("payroll cannot be nil")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	payroll.TenantID = tenantID
	err = r.db.WithContext(ctx).Create(payroll).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormPayrollRepo) GetByID(ctx context.Context, id uint) (*domain.Payroll, error) {
	if id == 0 {
		return nil, errors.New("invalid payroll id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}

	var payroll domain.Payroll
	err = r.db.WithContext(ctx).
		Preload("Employee.User").
		Preload("Items").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&payroll).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPayrollNotFound
		}
		return nil, err
	}
	return &payroll, nil

}

func (r *GormPayrollRepo) GetByEmployeeAndPeriod(ctx context.Context, employeID uint, periodStart, periodEnd time.Time) (*domain.Payroll, error) {
	if employeID == 0 {
		return nil, errors.New("invalid employeid")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var payroll domain.Payroll
	err = r.db.WithContext(ctx).
		Preload("Employee.User").
		Preload("Items").
		Where("tenant_id = ? AND employe_id = ? AND period_start = ? AND period_end =?", tenantID, employeID, periodStart, periodEnd).
		First(&payroll).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPayrollNotFound
		}
		return nil, err
	}
	return &payroll, nil
}

func (r *GormPayrollRepo) ListByEmployee(ctx context.Context, employeeID uint) ([]domain.Payroll, error) {
	if employeeID == 0 {
		return nil, errors.New("invalid employee id")
	}
	tenanID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var payrolls []domain.Payroll
	err = r.db.WithContext(ctx).
		Preload("Employee.User").
		Preload("Items").
		Where("tenant_id = ? AND employee_id = ?", tenanID, employeeID).
		Order("Period").
		Find(&payrolls).Error
	if err != nil {
		return nil, err
	}
	return payrolls, err
}

func (r *GormPayrollRepo) Update(ctx context.Context, payroll *domain.Payroll) error {
	if payroll == nil || payroll.ID == 0 {
		return errors.New("payroll cannot be nil or 0")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	existing, err := r.GetByID(ctx, payroll.ID)
	if err != nil {
		return err
	}
	payroll.TenantID = existing.TenantID
	err = r.db.WithContext(ctx).
		Model(&domain.Payroll{}).
		Where("id = ? AND tenant_id = ?", payroll.ID, tenantID).Updates(payroll).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormPayrollRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid payroll id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&domain.Payroll{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrPayrollNotFound
	}
	return nil
}
