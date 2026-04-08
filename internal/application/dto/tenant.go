package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateTenantRequest struct {
	Name   string `json:"name" binding:"required,min=1,max=255"`
	Active *bool  `json:"active"`
}

type TenantResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}
