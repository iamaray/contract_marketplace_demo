package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"contract_market_demo/backend/models"
	"contract_market_demo/backend/repos"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/google/uuid"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/joho/godotenv"
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
		// &models.PaymentsAccount{},
		&models.ContractListing{},
		&models.ContractHeader{},
		&models.ContractState{},
		&models.TransactionRecord{},
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

func GetListing(
	listingID uuid.UUID,
	listingRepo repos.ContractListingRepository) (*models.ContractListing, error) {
	listing, err := listingRepo.FindByID(listingID)
	if err != nil {
		return nil, err
	}

	return listing, nil
}

func CurrentUser(r *http.Request, userRepo repos.UserRepository) (*models.User, error) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		return nil, errors.New("unauthenticated")
	}
	sub := claims.Subject

	return userRepo.FindOrCreateByAuth("clerk", sub, "")
}

// func GetUser(userID uuid.UUID) {
// 	// retrieve the user from DB

// }

type ListingCreateRequest struct {
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
	stateRepo repos.ContractStateRepository,
	userRepo repos.UserRepository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			req := &ListingCreateRequest{}
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
				return
			}

			u, err := CurrentUser(r, userRepo)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			sellerID := u.ID
			// sellerID, err := uuid.Parse(req.SellerID)
			// if err != nil {
			// 	http.Error(w, "invalid seller_id: "+err.Error(), http.StatusBadRequest)
			// 	return
			// }

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
			req := &ListingGetRequest{}
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
				return
			}
			listingID := req.ListingID

			if listingID != "" {
				listingUUID, err := uuid.Parse(req.ListingID)
				if err != nil {
					http.Error(w, "invalid seller_id: "+err.Error(), http.StatusBadRequest)
					return
				}
				listing, err := GetListing(listingUUID, listingRepo)
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

			u, err := CurrentUser(r, userRepo)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if listing.SellerID != u.ID {
				http.Error(w, "forbidden", http.StatusForbidden)
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
			listing, err := GetListing(id, listingRepo)
			if err != nil {
				http.Error(w, "listing not found", http.StatusBadRequest)
				return
			}

			u, err := CurrentUser(r, userRepo)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if listing.SellerID != u.ID {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

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

func IssueFromListing(
	listing *models.ContractListing,
	issueQuantity int,
	listingRepo repos.ContractListingRepository) ([]*models.ContractHeader, []*models.ContractState, error) {

	if issueQuantity > int(listing.SupplyRemaining) {
		return nil, nil, errors.New("not enough supply")
	}

	listing.SupplyRemaining -= uint64(issueQuantity)
	err := listingRepo.Update(listing)
	if err != nil {
		// listing.SupplyRemaining += uint64(issueQuantity)
		return nil, nil, err
	}

	headers := make([]*models.ContractHeader, issueQuantity)
	states := make([]*models.ContractState, issueQuantity)
	for i := 0; i < issueQuantity; i++ {
		header := NewHeader(listing.ID)
		headers[i] = header

		state := &models.ContractState{
			HeaderID:       header.ID,
			LastPurchaseAt: time.Now(),
			OwnerID:        listing.SellerID,
			Status:         models.StatusListed,
		}
		states[i] = state
	}

	return headers, states, nil
}

func TransferOwnership(
	buyerID uuid.UUID,
	state *models.ContractState,
	userRepo repos.UserRepository,
	stateRepo repos.ContractStateRepository) (uuid.UUID, error) {
	_, err := userRepo.FindByID(buyerID)
	if err != nil {
		return uuid.Nil, err
	}

	prevOwner := state.OwnerID
	state.OwnerID = buyerID
	state.LastPurchaseAt = time.Now()
	state.Status = models.StatusOwned

	err = stateRepo.Update(state)
	if err != nil {
		return uuid.Nil, err
	}

	return prevOwner, nil
}

func SettleTransaction(record *models.TransactionRecord) (*models.TransactionRecord, error) {
	// Settle the transaction and emit a transaction record.
	// This function will be the main entrypoint to the
	// Stripe Connect logic.

	return record, nil
}

func Transact(
	listingID uuid.UUID,
	buyerID uuid.UUID,
	purchaseQuantity int,
	userRepo repos.UserRepository,
	transactionRepo repos.TransactionRepository,
	listingRepo repos.ContractListingRepository,
	headerRepo repos.ContractHeaderRepository,
	stateRepo repos.ContractStateRepository) (*models.TransactionRecord, error) {

	record := &models.TransactionRecord{
		ListingID:         listingID,
		SellerID:          uuid.Nil,
		BuyerID:           buyerID,
		PurchaseQuantity:  purchaseQuantity,
		TransactionStatus: models.StatusPending,
	}
	// handle this error by just logging it
	_ = transactionRepo.Create(record)

	listing, err := GetListing(listingID, listingRepo)
	if err != nil {
		record.TransactionStatus = models.StatusFailed
		// handle this error by just logging it
		_ = transactionRepo.Update(record)
		return record, err
	}
	record.SellerID = listing.SellerID
	// handle this error by just logging it
	_ = transactionRepo.Update(record)

	_, states, err := IssueFromListing(listing, purchaseQuantity, listingRepo)
	if err != nil {
		record.TransactionStatus = models.StatusFailed
		// handle this error by just logging it
		_ = transactionRepo.Update(record)
		return record, err
	}

	for i := 0; i < purchaseQuantity; i++ {
		_, err := TransferOwnership(buyerID, states[i], userRepo, stateRepo)
		if err != nil {
			record.TransactionStatus = models.StatusFailed
			// handle this error by just logging it
			_ = transactionRepo.Update(record)
			return record, err
		}
	}

	record, err = SettleTransaction(record)
	// handle this error by just logging it
	_ = transactionRepo.Update(record)

	if record.TransactionStatus == models.StatusFailed || err != nil {
		listing.SupplyRemaining += uint64(purchaseQuantity)
		err = listingRepo.Update(listing)
		if err != nil {
			return record, errors.New("error restoring supply after failed transaction")
		}
		return record, err
	}

	// commit the record to DB

	return record, err
}

type ContractPurchaseRequest struct {
	ListingID        string
	PurchaseQuantity int
}

func ContractPurchaseHandler(
	userRepo repos.UserRepository,
	transactionRepo repos.TransactionRepository,
	listingRepo repos.ContractListingRepository,
	headerRepo repos.ContractHeaderRepository,
	stateRepo repos.ContractStateRepository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			req := &ContractPurchaseRequest{}
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			}

			// buyerID, err := uuid.Parse(req.BuyerID)
			// if err != nil {
			// 	http.Error(w, "invalid buyer_id: "+err.Error(), http.StatusBadRequest)
			// 	return
			// }

			listingID, err := uuid.Parse(req.ListingID)
			if err != nil {
				http.Error(w, "invalid listing_id: "+err.Error(), http.StatusBadRequest)
				return
			}

			u, err := CurrentUser(r, userRepo)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
			}

			_, err = Transact(
				listingID,
				u.ID,
				req.PurchaseQuantity,
				userRepo,
				transactionRepo,
				listingRepo,
				headerRepo,
				stateRepo)

			if err != nil {
				http.Error(w, "transaction error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		if r.Method == http.MethodGet {
			// retrieve a transaction record.
			return
		}
		if r.Method == http.MethodPut {
			// update a transaction record
			return
		}
		if r.Method == http.MethodDelete {
			// delete a transaction record
			return
		}
	}
}

func MeHandler(userRepo repos.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := CurrentUser(r, userRepo)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		claims, _ := clerk.SessionClaimsFromContext(r.Context())
		resp := map[string]any{
			"user_id":  u.ID.String(),
			"subject":  claims.Subject,
			"provider": u.AuthProvider,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func main() {
	clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))
	_ = godotenv.Load()
	db := SetupDB()
	db.AutoMigrate()

	userRepo := repos.NewUserRepository(db.DB)
	transactionRepo := repos.NewTransactionRepository(db.DB)
	listingRepo := repos.NewContractListingRepository(db.DB)
	headerRepo := repos.NewContractHeaderRepository(db.DB)
	stateRepo := repos.NewContractStateRepository(db.DB)

	// _, _ = CreateDummySellers(10, db.DB)

	mux := http.NewServeMux()
	listings := http.HandlerFunc(HeaderListingHandler(
		listingRepo, headerRepo, stateRepo, userRepo))
	mux.Handle("/v1/listings", clerkhttp.WithHeaderAuthorization()(listings))

	contracts := http.HandlerFunc(ContractPurchaseHandler(
		userRepo, transactionRepo, listingRepo, headerRepo, stateRepo))
	mux.Handle("/v1/contracts", clerkhttp.RequireHeaderAuthorization()(contracts))

	mux.Handle("/v1/me", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(MeHandler(userRepo))))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	log.Println("listening on :8080")
	log.Println(srv.ListenAndServe())
}
