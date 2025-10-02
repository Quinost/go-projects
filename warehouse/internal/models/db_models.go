package models

import (
	"github.com/google/uuid"
)

type Item struct {
	Id          uuid.UUID
	Name        string
	Description string
}

type User struct {
	Id       uuid.UUID
	Username string
	Password string
}
