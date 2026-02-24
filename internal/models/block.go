package models

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"
)

type Block struct {
	ID uuid.UUID
	NoteID uuid.UUID
	Type string
	State json.RawMessage
}
