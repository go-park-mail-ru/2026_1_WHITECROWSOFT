package models

import (
	uuid "github.com/satori/go.uuid"
)

type Note struct {
	ID uuid.UUID
	UserID uuid.UUID
	IsSubnote bool
	NoteID uuid.UUID
}
