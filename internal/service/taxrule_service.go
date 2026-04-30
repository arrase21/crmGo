package service

import (
	"context"
	"errors"

	"github.com/arrase21/crm-users/internal/domain"
)

// ========================================
// Errors
// ========================================

var (
	ErrTaxRuleNotFound = errors.New("tax rule not found")
)

// ========================================
// TaxRule Service
// ========================================

type TaxRuleService struct {
	taxRuleRepo domain.TaxRuleRepo
}

func NewTaxRuleService(taxRuleRepo domain.TaxRuleRepo) *TaxRuleService {
	return &TaxRuleService{
		taxRuleRepo: taxRuleRepo,
	}
}

// TaxRuleRequest request para crear regla de retención
type TaxRuleRequest struct {
	Code        string  `json:"code" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	BasePercent float64 `json:"base_percent" validate:"required"`
	TaxCategory string  `json:"tax_category"`
	MinIncome   float64 `json:"min_income"`
	MaxIncome   float64 `json:"max_income"`
	Priority    int     `json:"priority"`
}

// CalculateWithholding calcula la retención en fuente
func (s *TaxRuleService) CalculateWithholding(baseIncome float64) (float64, error) {
	rule, err := s.taxRuleRepo.GetByIncome(context.Background(), baseIncome)
	if err != nil {
		return 0, err
	}

	// Aplicar retención solo si hay regla
	if rule.BasePercent > 0 {
		return baseIncome * (rule.BasePercent / 100), nil
	}

	return 0, nil
}

// Create crea una regla de retención
func (s *TaxRuleService) Create(ctx context.Context, req TaxRuleRequest) (*domain.TaxRule, error) {
	rule := &domain.TaxRule{
		Code:        req.Code,
		Name:        req.Name,
		BasePercent: req.BasePercent,
		TaxCategory: req.TaxCategory,
		MinIncome:   req.MinIncome,
		MaxIncome:   req.MaxIncome,
		Priority:    req.Priority,
		IsActive:    true,
	}

	if err := s.taxRuleRepo.Create(ctx, rule); err != nil {
		return nil, err
	}

	return rule, nil
}

// List lista las reglas de retención
func (s *TaxRuleService) List(ctx context.Context, page, limit int) ([]domain.TaxRule, int64, error) {
	return s.taxRuleRepo.List(ctx, page, limit)
}

// GetByIncome obtiene la regla aplicable para un ingreso
func (s *TaxRuleService) GetByIncome(ctx context.Context, income float64) (*domain.TaxRule, error) {
	return s.taxRuleRepo.GetByIncome(ctx, income)
}

// SeedDefaultTaxRules carga reglas de retención por defecto
func (s *TaxRuleService) SeedDefaultTaxRules(ctx context.Context) error {
	rules := []TaxRuleRequest{
		{
			Code:        "RET_FUENTE_0",
			Name:        "Sin retención (ingresos menores a 10 UVT)",
			BasePercent: 0,
			TaxCategory: "ingresos_trabajo",
			MinIncome:   0,
			MaxIncome:   1167000, // ~10 UVT
			Priority:    1,
		},
		{
			Code:        "RET_FUENTE_10",
			Name:        "Retención 10% (ingresos de 10 a 20 UVT)",
			BasePercent: 10,
			TaxCategory: "ingresos_trabajo",
			MinIncome:   1167000,
			MaxIncome:   2333400, // ~20 UVT
			Priority:    2,
		},
		{
			Code:        "RET_FUENTE_20",
			Name:        "Retención 20% (ingresos de 20 a 40 UVT)",
			BasePercent: 20,
			TaxCategory: "ingresos_trabajo",
			MinIncome:   2333400,
			MaxIncome:   4666800, // ~40 UVT
			Priority:    3,
		},
		{
			Code:        "RET_FUENTE_30",
			Name:        "Retención 30% (ingresos mayores a 40 UVT)",
			BasePercent: 30,
			TaxCategory: "ingresos_trabajo",
			MinIncome:   4666800,
			MaxIncome:   0, // Sin límite
			Priority:    4,
		},
	}

	for _, r := range rules {
		_, err := s.Create(ctx, r)
		if err != nil {
			// Ignorar errores de duplicado
			continue
		}
	}

	return nil
}
