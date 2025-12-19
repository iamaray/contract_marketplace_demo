package models

import (
	"github.com/google/uuid"
)

type Seller struct {
	ID uuid.UUID
}

type PaymentsAccount struct {
	Seller Seller
}
