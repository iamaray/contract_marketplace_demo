package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"contract_market_demo/backend/models"
	"contract_market_demo/backend/repos"

	"github.com/google/uuid"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB
}

func SetupDB() *Database {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	log.Println("Connected to database")
	return &Database{DB: db}
}

func (db *Database) AutoMigrate() {
	db.DB.AutoMigrate(
		&models.User{},
		&models.PaymentsAccount{},
		&models.ContractListing{},
		&models.ContractHeader{},
		&models.ContractState{},
	)
	log.Println("Database migration complete")
}

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
		HeaderID:       headerID,
		LastPurchaseAt: time.Now(),
		OwnerID:        ownerID,
		Status:         models.StatusOwned,
	}
}

func CreateListing(
	sellerID uuid.UUID,
	listPriceNanos int64,
	supplyLimit uint64,
	listingRepo repos.ContractListingRepository,
) (*models.ContractListing, error) {
	listing := NewListing(
		sellerID,
		listPriceNanos,
		supplyLimit,
	)

	err := listingRepo.Create(listing)
	if err != nil {
		return nil, err
	}

	return listing, nil
}

// func GetListing(listingID uuid.UUID) *models.ContractListing {

// }

type ListingCreateRequest struct {
	SellerID       string `json:"seller_id"`
	ListPriceNanos int64  `json:"list_price_nanos"`
	SupplyLimit    uint64 `json:"supply_limit"`
}

type ListingUpdateRequest struct {
	ListingID      string `json:"listing_id"`
	ListPriceNanos int64  `json:"list_price_nanos"`
	SupplyLimit    uint64 `json:"supply_limit"`
}

type ListingGetRequest struct {
	ListingID string `json:"listing_id"`
}

func HeaderListingHandler(
	listingRepo repos.ContractListingRepository,
	headerRepo repos.ContractHeaderRepository,
	stateRepo repos.ContractStateRepository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			req := &ListingCreateRequest{}
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
				return
			}
			sellerID, err := uuid.Parse(req.SellerID)
			if err != nil {
				http.Error(w, "invalid seller_id: "+err.Error(), http.StatusBadRequest)
				return
			}

			listing, err := CreateListing(
				sellerID,
				req.ListPriceNanos,
				req.SupplyLimit,
				listingRepo,
			)
			if err != nil {
				http.Error(w, "failed to create listing: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(listing)
			return
		}
		if r.Method == http.MethodGet {
			listingID := r.URL.Query().Get("listing_id")
			if listingID != "" {
				id, err := uuid.Parse(listingID)
				if err != nil {
					http.Error(w, "invalid listing_id: "+err.Error(), http.StatusBadRequest)
					return
				}
				listing, err := listingRepo.FindByID(id)
				if err != nil {
					http.Error(w, "listing not found: "+err.Error(), http.StatusNotFound)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(listing)
				return
			}
			listings, err := listingRepo.FindAll()
			if err != nil {
				http.Error(w, "failed to fetch listings: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(listings)
			return
		}
		if r.Method == http.MethodPut {
			// update a listing: expect JSON { listing_id, list_price_nanos, supply_limit }
			req := &ListingUpdateRequest{}
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
				return
			}
			id, err := uuid.Parse(req.ListingID)
			if err != nil {
				http.Error(w, "invalid listing_id: "+err.Error(), http.StatusBadRequest)
				return
			}
			listing, err := listingRepo.FindByID(id)
			if err != nil {
				http.Error(w, "listing not found: "+err.Error(), http.StatusNotFound)
				return
			}
			listing.ListPriceNanos = req.ListPriceNanos
			listing.SupplyLimit = req.SupplyLimit
			listing.UpdatedAt = time.Now()
			if err := listingRepo.Update(listing); err != nil {
				http.Error(w, "failed to update listing: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(listing)
			return
		}
		if r.Method == http.MethodDelete {
			listingID := r.URL.Query().Get("listing_id")
			if listingID == "" {
				req := &ListingGetRequest{}
				if err := json.NewDecoder(r.Body).Decode(req); err != nil {
					http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
					return
				}
				listingID = req.ListingID
			}
			id, err := uuid.Parse(listingID)
			if err != nil {
				http.Error(w, "invalid listing_id: "+err.Error(), http.StatusBadRequest)
				return
			}
			if err := listingRepo.Delete(id); err != nil {
				http.Error(w, "failed to delete listing: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}

func ContractPurchaseHandler(
	listingRepo repos.ContractListingRepository,
	headerRepo repos.ContractHeaderRepository,
	stateRepo repos.ContractStateRepository) http.HandlerFunc {

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
	db := SetupDB()
	db.AutoMigrate()

	listing_repo := repos.NewContractListingRepository(db.DB)
	header_repo := repos.NewContractHeaderRepository(db.DB)
	state_repo := repos.NewContractStateRepository(db.DB)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/listings", HeaderListingHandler(
		listing_repo, header_repo, state_repo))
	mux.HandleFunc("/v1/contracts", ContractPurchaseHandler(
		listing_repo, header_repo, state_repo))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	log.Println("listening on :8080")
	log.Println(srv.ListenAndServe())
}
