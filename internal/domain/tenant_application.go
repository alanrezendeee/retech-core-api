package domain

import "github.com/google/uuid"

type TenantApplication struct {
	TenantID      uuid.UUID
	ApplicationID uuid.UUID
}
