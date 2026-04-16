package repository

import (
	"context"
	"errors"
	"time"

	"github.com/arrase21/crm-users/internal/domain"
	"gorm.io/gorm"
)

type GormEmployeeContractRepo struct {
	db *gorm.DB
}

func NewGormEmployeeContractRepository(db *gorm.DB) domain.EmployeeContractRepo {
	return &GormEmployeeContractRepo{
		db: db,
	}
}

func (r *GormEmployeeContractRepo) Create(ctx context.Context, contract *domain.EmployeeContract) error {
	if contract == nil {
		return errors.New("contract cannot be nil")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	contract.TenantID = tenantID
	err = r.db.WithContext(ctx).Create(contract).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormEmployeeContractRepo) GetByID(ctx context.Context, id uint) (*domain.EmployeeContract, error) {
	if id == 0 {
		return nil, errors.New("invalid contract id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}

	var contract domain.EmployeeContract
	err = r.db.WithContext(ctx).
		Preload("Employee.User").
		Preload("ContractType").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&contract).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEmployeeContractNotFound
		}
		return nil, err
	}
	return &contract, nil
}

func (r *GormEmployeeContractRepo) GetActiveByEmployee(ctx context.Context, employeeID uint) (*domain.EmployeeContract, error) {
	if employeeID == 0 {
		return nil, errors.New("invalid employee id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}

	var contract domain.EmployeeContract
	err = r.db.WithContext(ctx).
		Preload("Employee.User").
		Preload("ContractType").
		Where("tenant_id = ? AND employee_id = ? AND is_active = ?", tenantID, employeeID, true).
		First(&contract).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEmployeeContractNotFound
		}
		return nil, err
	}
	return &contract, nil
}

func (r *GormEmployeeContractRepo) ListByEmployee(ctx context.Context, employeeID uint) ([]domain.EmployeeContract, error) {
	if employeeID == 0 {
		return nil, errors.New("invalid employee id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var contracts []domain.EmployeeContract
	err = r.db.WithContext(ctx).
		Preload("Employee.User").
		Preload("ContractType").
		Where("tenant_id = ? AND employee_id = ?", tenantID, employeeID).
		Order("start_date DESC").
		Find(&contracts).Error
	if err != nil {
		return nil, err
	}
	return contracts, nil
}

func (r *GormEmployeeContractRepo) Update(ctx context.Context, contract *domain.EmployeeContract) error {
	if contract == nil || contract.ID == 0 {
		return errors.New("contract cannot be nil or 0")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	existing, err := r.GetByID(ctx, contract.ID)
	if err != nil {
		return err
	}
	contract.TenantID = existing.TenantID
	err = r.db.WithContext(ctx).
		Model(&domain.EmployeeContract{}).
		Where("id = ? AND tenant_id = ?", contract.ID, tenantID).Updates(contract).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormEmployeeContractRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid contract id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&domain.EmployeeContract{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrEmployeeContractNotFound
	}
	return nil
}

// Helper to parse contract dates
func parseContractDates(contract *domain.EmployeeContract, periodStart, periodEnd time.Time) bool {
	// Check if contract is active during the period
	if contract.StartDate.After(periodEnd) {
		return false
	}
	if contract.EndDate != nil && contract.EndDate.Before(periodStart) {
		return false
	}
	return true
}
