package domain

import "github.com/google/uuid"

type UserTenant struct {
	UserID   uuid.UUID
	TenantID uuid.UUID
}
