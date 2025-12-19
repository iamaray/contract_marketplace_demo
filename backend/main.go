package main

import (
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
	db := SetupDB()
	db.AutoMigrate()

	repos.NewContractListingRepository(db.DB)
	repos.NewContractHeaderRepository(db.DB)
	repos.NewContractStateRepository(db.DB)

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
