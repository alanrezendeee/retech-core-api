package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/theretech/retech-core-api/internal/application/dto"
	"github.com/theretech/retech-core-api/internal/domain"
)

type Application struct {
	apps domain.ApplicationRepository
	ta   domain.TenantApplicationRepository
}

func NewApplication(apps domain.ApplicationRepository, ta domain.TenantApplicationRepository) *Application {
	return &Application{apps: apps, ta: ta}
}

func (u *Application) Create(ctx context.Context, req dto.CreateApplicationRequest) (*dto.ApplicationResponse, error) {
	active := true
	if req.Active != nil {
		active = *req.Active
	}
	now := time.Now().UTC()
	a := &domain.Application{
		ID:        uuid.New(),
		Name:      req.Name,
		Active:    active,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := u.apps.Create(ctx, a); err != nil {
		return nil, err
	}
	return toApplicationResponse(a), nil
}

// ListForTenant retorna apenas aplicações ativas vinculadas ao tenant.
func (u *Application) ListForTenant(ctx context.Context, tenantID uuid.UUID) ([]dto.ApplicationResponse, error) {
	list, err := u.apps.ListByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ApplicationResponse, 0, len(list))
	for i := range list {
		if !list[i].Active {
			continue
		}
		out = append(out, *toApplicationResponse(&list[i]))
	}
	return out, nil
}

func toApplicationResponse(a *domain.Application) *dto.ApplicationResponse {
	return &dto.ApplicationResponse{
		ID:        a.ID,
		Name:      a.Name,
		Active:    a.Active,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
