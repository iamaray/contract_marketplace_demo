package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	Email string `gorm:"index"`

	AuthProvider string `gorm:"size:32;not null;index:idx_auth,unique"`
	AuthSubject  string `gorm:"size:191;not null;index:idx_auth,unique"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// type PaymentsAccount struct {
// 	User User
// }
