package models

import (
	"time"

	"github.com/google/uuid"
)

type BlockState struct {
	ID         uuid.UUID
	BlockID    uuid.UUID
	Formatting string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
