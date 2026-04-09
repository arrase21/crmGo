package repository

import (
	"context"
	"errors"

	"github.com/arrase21/crm-users/internal/domain"
	"gorm.io/gorm"
)

type GormPayrollConceptRepo struct {
	db *gorm.DB
}

func NewGormPayrollConceptRepository(db *gorm.DB) domain.PayrollConceptRepo {
	return &GormPayrollConceptRepo{
		db: db,
	}
}

func (r *GormPayrollConceptRepo) Create(ctx context.Context, concept *domain.PayrollConcept) error {
	if concept == nil {
		return errors.New("concept cannot be nil")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	concept.TenantID = tenantID
	err = r.db.WithContext(ctx).Create(concept).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *GormPayrollConceptRepo) GetByID(ctx context.Context, id uint) (*domain.PayrollConcept, error) {
	if id == 0 {
		return nil, errors.New("invald payroll concept id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var concept domain.PayrollConcept
	err = r.db.WithContext(ctx).Where("tenant_id =? AND id =?", tenantID, id).First(&concept).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrConceptNotFound
		}
		return nil, err
	}
	return &concept, nil
}

func (r *GormPayrollConceptRepo) GetByCode(ctx context.Context, code string) (*domain.PayrollConcept, error) {
	if code == "" {
		return nil, errors.New("invalid code")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var concept domain.PayrollConcept
	err = r.db.WithContext(ctx).
		Where("code = ? AND tenant_id = ?", code, tenantID).
		First(&concept).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrConceptNotFound
		}
		return nil, err
	}
	return &concept, nil
}

func (r *GormPayrollConceptRepo) GetActiveConcepts(ctx context.Context) ([]domain.PayrollConcept, error) {
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, err
	}
	var concepts []domain.PayrollConcept
	err = r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Order("code").
		Find(&concepts).Error
	if err != nil {
		return nil, err
	}
	return concepts, nil
}

func (r *GormPayrollConceptRepo) List(ctx context.Context, page, limit int) ([]domain.PayrollConcept, int64, error) {
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var concepts []domain.PayrollConcept
	err = r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Offset(offset).
		Limit(limit).
		Order("code").
		Find(&concepts).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	r.db.Model(&domain.PayrollConcept{}).
		Where("tenant_id = ?", tenantID).
		Count(&total)

	return concepts, total, nil
}

func (r *GormPayrollConceptRepo) Update(ctx context.Context, concept *domain.PayrollConcept) error {
	if concept == nil || concept.ID == 0 {
		return errors.New("concept cannot be nil or with zero id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	err = r.db.WithContext(ctx).
		Model(&domain.PayrollConcept{}).
		Where("id = ? AND tenant_id = ?", concept.ID, tenantID).
		Updates(map[string]interface{}{
			"name":          concept.Name,
			"type":          concept.Type,
			"description":   concept.Description,
			"percentage":    concept.Percentage,
			"employee_part": concept.EmployeePart,
			"employer_part": concept.EmployerPart,
			"is_mandatory":  concept.IsMandatory,
			"is_active":     concept.IsActive,
		}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormPayrollConceptRepo) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid payroll concept id")
	}
	tenantID, err := tenantFromctx(ctx)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&domain.PayrollConcept{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrConceptNotFound
	}
	return nil
}
