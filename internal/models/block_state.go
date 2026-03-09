package models

import (
	"time"

	"github.com/google/uuid"
)

type BlockState struct {
	ID         uuid.UUID `json:"id" db:"id"`
	BlockID    uuid.UUID `json:"block_id" db:"block_id"`
	Formatting string    `json:"formatting" db:"formatting"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
