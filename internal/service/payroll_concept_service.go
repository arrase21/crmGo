package service

import (
	"context"
	"errors"

	"github.com/arrase21/crm-users/internal/domain"
)

type PayrollConceptService struct {
	conceptRepo domain.PayrollConceptRepo
}

func NewPayrollConceptService(repo domain.PayrollConceptRepo) *PayrollConceptService {
	return &PayrollConceptService{
		conceptRepo: repo,
	}
}

func (s *PayrollConceptService) Create(ctx context.Context, concept *domain.PayrollConcept) error {
	if concept == nil {
		return errors.New("concept cannot be nil")
	}
	return s.conceptRepo.Create(ctx, concept)
}

func (s *PayrollConceptService) GetByID(ctx context.Context, id uint) (*domain.PayrollConcept, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	return s.conceptRepo.GetByID(ctx, id)
}

func (s *PayrollConceptService) GetByCode(ctx context.Context, code string) (*domain.PayrollConcept, error) {
	if code == "" {
		return nil, errors.New("invalid code")
	}
	return s.conceptRepo.GetByCode(ctx, code)
}

func (s *PayrollConceptService) GetActiveConcepts(ctx context.Context) ([]domain.PayrollConcept, error) {
	return s.conceptRepo.GetActiveConcepts(ctx)
}

func (s *PayrollConceptService) List(ctx context.Context, page, limit int) ([]domain.PayrollConcept, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	return s.conceptRepo.List(ctx, page, limit)
}

func (s *PayrollConceptService) Update(ctx context.Context, concept *domain.PayrollConcept) error {
	if concept == nil {
		return errors.New("concept cannot be nil")
	}
	if concept.ID == 0 {
		return errors.New("invalid concept id")
	}
	return s.conceptRepo.Update(ctx, concept)
}

func (s *PayrollConceptService) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid id")
	}
	return s.conceptRepo.Delete(ctx, id)
}

func (s *PayrollConceptService) SeedDefaultConcepts(ctx context.Context) error {
	defaults := domain.DefaultPayrollConcepts()

	for _, concept := range defaults {
		// Verificar si ya existe
		existing, err := s.conceptRepo.GetByCode(ctx, concept.Code)
		if err == nil && existing != nil {
			continue // Ya existe, skip
		}

		if err := s.conceptRepo.Create(ctx, &concept); err != nil {
			return err
		}
	}

	return nil
}
