package persistence

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserTenantRepository struct {
	db *gorm.DB
}

func NewUserTenantRepository(db *gorm.DB) *UserTenantRepository {
	return &UserTenantRepository{db: db}
}

func (r *UserTenantRepository) exists(ctx context.Context, userID, tenantID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&userTenantModel{}).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserTenantRepository) Link(ctx context.Context, userID, tenantID uuid.UUID) error {
	ok, err := r.exists(ctx, userID, tenantID)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	m := userTenantModel{
		UserID:    userID,
		TenantID:  tenantID,
		CreatedAt: time.Now().UTC(),
	}
	return r.db.WithContext(ctx).Create(&m).Error
}
