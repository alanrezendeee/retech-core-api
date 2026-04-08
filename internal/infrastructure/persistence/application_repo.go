package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/theretech/retech-core-api/internal/domain"
)

type ApplicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) Create(ctx context.Context, a *domain.Application) error {
	m := applicationModel{
		ID:        a.ID,
		Name:      a.Name,
		Active:    a.Active,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
	err := r.db.WithContext(ctx).Create(&m).Error
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		return err
	}
	return nil
}

func (r *ApplicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Application, error) {
	var m applicationModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &domain.Application{
		ID:        m.ID,
		Name:      m.Name,
		Active:    m.Active,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

func (r *ApplicationRepository) ListByTenantID(ctx context.Context, tenantID uuid.UUID) ([]domain.Application, error) {
	var rows []applicationModel
	err := r.db.WithContext(ctx).
		Table(applicationModel{}.TableName()+" AS a").
		Select("a.*").
		Joins("INNER JOIN tenant_applications ta ON ta.application_id = a.id").
		Where("ta.tenant_id = ?", tenantID).
		Order("a.name ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.Application, 0, len(rows))
	for i := range rows {
		out = append(out, domain.Application{
			ID:        rows[i].ID,
			Name:      rows[i].Name,
			Active:    rows[i].Active,
			CreatedAt: rows[i].CreatedAt,
			UpdatedAt: rows[i].UpdatedAt,
		})
	}
	return out, nil
}
