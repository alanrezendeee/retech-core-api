package domain

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID        uuid.UUID
	Name      string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
