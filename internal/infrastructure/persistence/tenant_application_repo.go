package persistence

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantApplicationRepository struct {
	db *gorm.DB
}

func NewTenantApplicationRepository(db *gorm.DB) *TenantApplicationRepository {
	return &TenantApplicationRepository{db: db}
}

func (r *TenantApplicationRepository) IsLinked(ctx context.Context, tenantID, applicationID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&tenantApplicationModel{}).
		Where("tenant_id = ? AND application_id = ?", tenantID, applicationID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *TenantApplicationRepository) Link(ctx context.Context, tenantID, applicationID uuid.UUID) error {
	ok, err := r.IsLinked(ctx, tenantID, applicationID)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	m := tenantApplicationModel{
		TenantID:      tenantID,
		ApplicationID: applicationID,
		CreatedAt:     time.Now().UTC(),
	}
	return r.db.WithContext(ctx).Create(&m).Error
}
