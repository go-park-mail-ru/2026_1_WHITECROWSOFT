package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Block struct {
	ID uuid.UUID `json:"id"`
	NoteID uuid.UUID `json:"note_id"`
	Type string `json:"type"`
	State json.RawMessage `json:"state"`
}
