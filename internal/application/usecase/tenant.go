package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/theretech/retech-core-api/internal/application/dto"
	"github.com/theretech/retech-core-api/internal/domain"
)

type Tenant struct {
	repo domain.TenantRepository
}

func NewTenant(repo domain.TenantRepository) *Tenant {
	return &Tenant{repo: repo}
}

func (u *Tenant) Create(ctx context.Context, req dto.CreateTenantRequest) (*dto.TenantResponse, error) {
	active := true
	if req.Active != nil {
		active = *req.Active
	}
	t := &domain.Tenant{
		ID:        uuid.New(),
		Name:      req.Name,
		Active:    active,
		CreatedAt: time.Now().UTC(),
	}
	if err := u.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return toTenantResponse(t), nil
}

func (u *Tenant) GetByID(ctx context.Context, id uuid.UUID) (*dto.TenantResponse, error) {
	t, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toTenantResponse(t), nil
}

func toTenantResponse(t *domain.Tenant) *dto.TenantResponse {
	return &dto.TenantResponse{
		ID:        t.ID,
		Name:      t.Name,
		Active:    t.Active,
		CreatedAt: t.CreatedAt,
	}
}
