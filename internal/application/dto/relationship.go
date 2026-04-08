package dto

import "github.com/google/uuid"

type LinkApplicationRequest struct {
	ApplicationID uuid.UUID `json:"application_id" binding:"required"`
}

type LinkUserRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}
