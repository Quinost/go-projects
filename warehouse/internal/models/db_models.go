package models

import (
	"github.com/google/uuid"
)

type Item struct {
	Id          uuid.UUID
	Name        string `validator:"required"`
	Description string `validator:"required"`
}

type User struct {
	Id       uuid.UUID
	Username string `validator:"required"`
	Password string `validator:"required"`
}
