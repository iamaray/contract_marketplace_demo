package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID
}

type PaymentsAccount struct {
	User User
}
