package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           uuid.UUID
	Username     string
	Password     []byte
	TokenVersion int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
