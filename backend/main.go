package main

import (
	"log"
	"net/http"
	"time"

	"contract_market_demo/backend/models"

	"github.com/google/uuid"
)

func NewListing(
	sellerID uuid.UUID,
	listPriceNanos int64,
	supplyLimit uint64) *models.ContractListing {

	return &models.ContractListing{
		ID:              uuid.New(),
		SellerID:        sellerID,
		ListPriceNanos:  listPriceNanos,
		SupplyLimit:     supplyLimit,
		SupplyRemaining: supplyLimit,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func NewHeader(
	listingID uuid.UUID,
) *models.ContractHeader {
	return &models.ContractHeader{
		ID:        uuid.New(),
		ListingID: listingID,
		CreatedAt: time.Now(),
	}
}

func NewState(
	headerID uuid.UUID,
	ownerID uuid.UUID,
) *models.ContractState {
	return &models.ContractState{
		HeaderID: headerID,
		LastPurchaseAt: time.Now(),
		OwnerID: ownerID,
		Status: models.StatusOwned,
	}
}

// func GetListing(listingID uuid.UUID) *models.ContractListing {

// }

func HeaderListingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			return
		}
		if r.Method == http.MethodGet {
			return
		}
		if r.Method == http.MethodPut {
			return
		}
		if r.Method == http.MethodDelete {
			return
		}
	}
}

func ContractPurchaseHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// initiate purchase
			return
		}
		if r.Method == http.MethodGet {
			return
		}
		if r.Method == http.MethodPut {
			return
		}
		if r.Method == http.MethodDelete {
			return
		}
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/listings", HeaderListingHandler())
	mux.HandleFunc("/v1/contracts", ContractPurchaseHandler())

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	log.Println("listening on :8080")
	log.Println(srv.ListenAndServe())
}
