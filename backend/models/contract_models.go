package models

import (
	"time"

	"github.com/google/uuid"
)

type ContractStatus uint8

const (
	StatusDraft ContractStatus = iota
	StatusListed
	StatusMatched
	StatusOwned
	StatusUnlocked
	StatusExpiryReached
)

type TransactionStatus uint8

const (
	StatusPending TransactionStatus = iota
	StatusRequiresPayment
	StatusPaid
	StatusFulfilled
	StatusExpired
	StatusFailed
)

type ContractListing struct {
	ID              uuid.UUID
	SellerID        uuid.UUID
	ListPriceNanos  int64
	SupplyLimit     uint64
	SupplyRemaining uint64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ContractHeader struct {
	ID        uuid.UUID
	ListingID uuid.UUID
	CreatedAt time.Time
}

type ContractState struct {
	HeaderID       uuid.UUID
	LastPurchaseAt time.Time
	OwnerID        uuid.UUID
	Status         ContractStatus
}

type TransactionRecord struct {
	ID				  uuid.UUID
	InitiatedAt       time.Time
	ListingID         uuid.UUID
	SellerID          uuid.UUID
	BuyerID           uuid.UUID
	PurchaseQuantity  int64
	PurchaseCents     uint64
	TransactionStatus TransactionStatus

	StripeCheckoutSessonID string
	StripePaymentIntentID  string
	Currency               string
	PlatformFeeCents       int64
	FulfilledAt            *time.Time
	IsFulfilled            bool
	StripeEventLastID      string
}
