package repository

import (
	"context"
	"errors"

	"github.com/arrase21/crm-users/internal/domain"
	"gorm.io/gorm"
)

type GormEmployeeRepo struct {
	db *gorm.DB
}

func NewGormEmployeeRepository(db *gorm.DB) domain.EmployeeRepo {
	return &GormEmployeeRepo{
		db: db,
	}
}

func (r *GormEmployeeRepo) Create(ctx context.Context, emp *domain.Employee) error {
	if emp == nil {
		return errors.New("employee cannot be nil")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	emp.TenantID = tenantID
	err = r.db.WithContext(ctx).Create(emp).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *GormEmployeeRepo) GetByID(ctx context.Context, id uint) (*domain.Employee, error) {
	if id == 0 {
		return nil, errors.New("invalid employee id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var employee domain.Employee
	err = r.db.WithContext(ctx).
		Preload("User").Preload("Department").Preload("Position").
		Preload("Contracts").Where("tenant_id = ? AND id =?", tenantID, id).First(&employee).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEmployeeNotFound
		}
		return nil, err
	}
	return &employee, nil
}

func (r *GormEmployeeRepo) GetByUserID(ctx context.Context, userID uint) (*domain.Employee, error) {
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var employee domain.Employee
	err = r.db.WithContext(ctx).
		Preload("User").
		Preload("Department").
		Preload("Position").
		Preload("Contracts").
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		First(&employee).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEmployeeNotFound
		}
		return nil, err
	}
	return &employee, nil
}

func (r *GormEmployeeRepo) List(ctx context.Context, page, limit int) ([]domain.Employee, int64, error) {
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	var employees []domain.Employee
	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.Employee{}).Where("tenant_id = ?", tenantID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Department").
		Preload("Position").
		Where("tenant_id = ?", tenantID).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&employees).Error; err != nil {
		return nil, 0, err
	}
	return employees, total, nil
}

func (r *GormEmployeeRepo) Update(ctx context.Context, emp *domain.Employee) error {
	if emp == nil {
		return errors.New("employee cannot be nil")
	}
	if emp.ID == 0 {
		return errors.New("invalid employee id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	emp.TenantID = tenantID
	err = r.db.WithContext(ctx).Save(emp).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormEmployeeRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid employee id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&domain.Employee{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrEmployeeNotFound
	}
	return nil
}

// ListActive retorna solo empleados activos con sus relaciones
func (r *GormEmployeeRepo) ListActive(ctx context.Context, page, limit int) ([]domain.Employee, int64, error) {
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 1000 {
		limit = 100 // Límite más alto para batch
	}
	offset := (page - 1) * limit

	var employees []domain.Employee
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Employee{}).
		Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Contracts", "is_active = ?", true). // Solo contratos activos
		Preload("Contracts.ContractType").
		Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Order("id ASC").
		Offset(offset).
		Limit(limit).
		Find(&employees).Error; err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}
