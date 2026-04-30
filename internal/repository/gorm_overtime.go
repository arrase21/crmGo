package repository

import (
	"context"
	"errors"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
	"gorm.io/gorm"
)

type GormOvertimeRepo struct {
	db *gorm.DB
}

func NewGormOvertimeRepository(db *gorm.DB) domain.OvertimeRepo {
	return &GormOvertimeRepo{db: db}
}

func (r *GormOvertimeRepo) Create(ctx context.Context, overtime *domain.Overtime) error {
	if overtime == nil {
		return errors.New("overtime cannot be nil")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	overtime.TenantID = tenantID
	return r.db.WithContext(ctx).Create(overtime).Error
}

func (r *GormOvertimeRepo) GetByID(ctx context.Context, id uint) (*domain.Overtime, error) {
	if id == 0 {
		return nil, errors.New("invalid overtime id")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var overtime domain.Overtime
	err = r.db.WithContext(ctx).
		Preload("Employee.User").
		Where("tenant_id = ? AND id = ?", tenantID, id).First(&overtime).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("overtime not found")
		}
		return nil, err
	}
	return &overtime, nil
}

func (r *GormOvertimeRepo) GetByEmployeeAndPeriod(ctx context.Context, employeeID uint, start, end time.Time) ([]domain.Overtime, error) {
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var overtimes []domain.Overtime
	err = r.db.WithContext(ctx).
		Where("tenant_id = ? AND employee_id = ? AND date >= ? AND date <= ?",
			tenantID, employeeID, start, end).
		Order("date ASC").
		Find(&overtimes).Error
	return overtimes, err
}

func (r *GormOvertimeRepo) GetByPayrollID(ctx context.Context, payrollID uint) ([]domain.Overtime, error) {
	var overtimes []domain.Overtime
	err := r.db.WithContext(ctx).
		Where("payroll_id = ?", payrollID).
		Find(&overtimes).Error
	return overtimes, err
}

func (r *GormOvertimeRepo) List(ctx context.Context, employeeID uint, page, limit int) ([]domain.Overtime, int64, error) {
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

	var overtimes []domain.Overtime
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Overtime{}).Where("tenant_id = ?", tenantID)
	if employeeID != 0 {
		query = query.Where("employee_id = ?", employeeID)
	}

	query.Count(&total)

	err = query.Preload("Employee.User").
		Order("date DESC").
		Offset(offset).
		Limit(limit).
		Find(&overtimes).Error

	return overtimes, total, err
}

func (r *GormOvertimeRepo) Update(ctx context.Context, overtime *domain.Overtime) error {
	if overtime == nil || overtime.ID == 0 {
		return errors.New("overtime cannot be nil or without id")
	}
	return r.db.WithContext(ctx).Save(overtime).Error
}

func (r *GormOvertimeRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid overtime id")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&domain.Overtime{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("overtime not found")
	}
	return nil
}
