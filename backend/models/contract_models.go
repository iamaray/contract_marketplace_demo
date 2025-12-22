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
	StatusExpired
)

type TransactionStatus uint8

const (
	StatusPending TransactionStatus = iota
	StatusSuccess
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
	ListingID         uuid.UUID
	SellerID          uuid.UUID
	BuyerID           uuid.UUID
	PurchaseQuantity  int
	PurchasePrice     uint64
	TransactionStatus TransactionStatus
}
