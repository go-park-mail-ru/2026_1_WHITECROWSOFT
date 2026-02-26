package models

import (
	"github.com/google/uuid"
)

type Note struct {
	ID       uuid.UUID   `json:"id"`
	UserID   uuid.UUID   `json:"user_id"`
	Title    string      `json:"title"`
	ParentID *uuid.UUID  `json:"parent_id,omitempty"`
	Blocks   []uuid.UUID `json:"blocks"`
}
