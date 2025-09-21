package models

import (
	"github.com/google/uuid"
)

type Item struct {
	Id          uuid.UUID
	Name        string
	Description string
}
