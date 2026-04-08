package domain

import (
	"context"

	"github.com/google/uuid"
)

type TenantRepository interface {
	Create(ctx context.Context, t *Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
}

type ApplicationRepository interface {
	Create(ctx context.Context, a *Application) error
	GetByID(ctx context.Context, id uuid.UUID) (*Application, error)
	ListByTenantID(ctx context.Context, tenantID uuid.UUID) ([]Application, error)
}

type TenantApplicationRepository interface {
	Link(ctx context.Context, tenantID, applicationID uuid.UUID) error
	IsLinked(ctx context.Context, tenantID, applicationID uuid.UUID) (bool, error)
}

type UserTenantRepository interface {
	Link(ctx context.Context, userID, tenantID uuid.UUID) error
}
