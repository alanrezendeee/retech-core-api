package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/theretech/retech-core-api/internal/domain"
)

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(ctx context.Context, t *domain.Tenant) error {
	m := tenantModel{
		ID:        t.ID,
		Name:      t.Name,
		Active:    t.Active,
		CreatedAt: t.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *TenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	var m tenantModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &domain.Tenant{
		ID:        m.ID,
		Name:      m.Name,
		Active:    m.Active,
		CreatedAt: m.CreatedAt,
	}, nil
}
