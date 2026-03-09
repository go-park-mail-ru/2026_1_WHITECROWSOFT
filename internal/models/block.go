package models

import (
	"time"

	"github.com/google/uuid"
)

type Block struct {
	ID          uuid.UUID `json:"id" db:"id"`
	NoteID      uuid.UUID `json:"note_id" db:"note_id"`
	BlockTypeID int       `json:"block_type_id" db:"block_type_id"`
	Position    int       `json:"position" db:"position"`
	Content     string    `json:"content" db:"content"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
