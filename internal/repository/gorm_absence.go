package repository

import (
	"context"
	"errors"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
	"gorm.io/gorm"
)

type GormAbsenceRepo struct {
	db *gorm.DB
}

func NewGormAbsenceRepository(db *gorm.DB) domain.AbsenceRepo {
	return &GormAbsenceRepo{db: db}
}

func (r *GormAbsenceRepo) Create(ctx context.Context, absence *domain.Absence) error {
	if absence == nil {
		return errors.New("absence cannot be nil")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	absence.TenantID = tenantID
	return r.db.WithContext(ctx).Create(absence).Error
}

func (r *GormAbsenceRepo) GetByID(ctx context.Context, id uint) (*domain.Absence, error) {
	if id == 0 {
		return nil, errors.New("invalid absence id")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var absence domain.Absence
	err = r.db.WithContext(ctx).
		Preload("Employee.User").
		Where("tenant_id = ? AND id = ?", tenantID, id).First(&absence).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("absence not found")
		}
		return nil, err
	}
	return &absence, nil
}

func (r *GormAbsenceRepo) GetByEmployeeAndPeriod(ctx context.Context, employeeID uint, start, end time.Time) ([]domain.Absence, error) {
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var absences []domain.Absence
	err = r.db.WithContext(ctx).
		Where("tenant_id = ? AND employee_id = ? AND start_date <= ? AND end_date >= ?",
			tenantID, employeeID, end, start).
		Order("start_date ASC").
		Find(&absences).Error
	return absences, err
}

func (r *GormAbsenceRepo) GetByPayrollID(ctx context.Context, payrollID uint) ([]domain.Absence, error) {
	var absences []domain.Absence
	err := r.db.WithContext(ctx).
		Where("payroll_id = ?", payrollID).
		Find(&absences).Error
	return absences, err
}

func (r *GormAbsenceRepo) List(ctx context.Context, employeeID uint, page, limit int) ([]domain.Absence, int64, error) {
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

	var absences []domain.Absence
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Absence{}).Where("tenant_id = ?", tenantID)
	if employeeID != 0 {
		query = query.Where("employee_id = ?", employeeID)
	}

	query.Count(&total)

	err = query.Preload("Employee.User").
		Order("start_date DESC").
		Offset(offset).
		Limit(limit).
		Find(&absences).Error

	return absences, total, err
}

func (r *GormAbsenceRepo) Update(ctx context.Context, absence *domain.Absence) error {
	if absence == nil || absence.ID == 0 {
		return errors.New("absence cannot be nil or without id")
	}
	return r.db.WithContext(ctx).Save(absence).Error
}

func (r *GormAbsenceRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid absence id")
	}
	tenantID, err := tenantFromCtx(ctx)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&domain.Absence{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("absence not found")
	}
	return nil
}
