package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/theretech/retech-core-api/internal/domain"
)

type Relationship struct {
	tenants domain.TenantRepository
	apps    domain.ApplicationRepository
	ta      domain.TenantApplicationRepository
	ut      domain.UserTenantRepository
}

func NewRelationship(
	tenants domain.TenantRepository,
	apps domain.ApplicationRepository,
	ta domain.TenantApplicationRepository,
	ut domain.UserTenantRepository,
) *Relationship {
	return &Relationship{
		tenants: tenants,
		apps:    apps,
		ta:      ta,
		ut:      ut,
	}
}

func (u *Relationship) LinkApplication(ctx context.Context, pathTenantID, headerTenantID, applicationID uuid.UUID) error {
	if pathTenantID != headerTenantID {
		return domain.ErrForbidden
	}
	tenant, err := u.tenants.GetByID(ctx, pathTenantID)
	if err != nil {
		return err
	}
	if !tenant.Active {
		return domain.ErrTenantInactive
	}
	app, err := u.apps.GetByID(ctx, applicationID)
	if err != nil {
		return err
	}
	if !app.Active {
		return domain.ErrApplicationInactive
	}
	return u.ta.Link(ctx, pathTenantID, applicationID)
}

func (u *Relationship) LinkUser(ctx context.Context, pathTenantID, headerTenantID, userID uuid.UUID) error {
	if pathTenantID != headerTenantID {
		return domain.ErrForbidden
	}
	tenant, err := u.tenants.GetByID(ctx, pathTenantID)
	if err != nil {
		return err
	}
	if !tenant.Active {
		return domain.ErrTenantInactive
	}
	return u.ut.Link(ctx, userID, pathTenantID)
}
