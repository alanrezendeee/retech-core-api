package persistence

import (
	"time"

	"github.com/google/uuid"
)

type tenantModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"size:255;not null"`
	Active    bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"not null"`
}

func (tenantModel) TableName() string { return "tenants" }

type applicationModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"size:128;not null;uniqueIndex"`
	Active    bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (applicationModel) TableName() string { return "applications" }

type tenantApplicationModel struct {
	TenantID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	ApplicationID uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt     time.Time `gorm:"not null"`
}

func (tenantApplicationModel) TableName() string { return "tenant_applications" }

type userTenantModel struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"not null"`
}

func (userTenantModel) TableName() string { return "user_tenants" }
