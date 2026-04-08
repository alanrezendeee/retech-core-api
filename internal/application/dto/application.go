package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateApplicationRequest struct {
	Name   string `json:"name" binding:"required,min=1,max=128"`
	Active *bool  `json:"active"`
}

type ApplicationResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
