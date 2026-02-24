package models

import (
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID       uuid.UUID
	Username string
	Password []byte
}
