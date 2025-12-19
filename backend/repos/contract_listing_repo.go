package repos

import (
	"contract_market_demo/backend/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrListingNotFound = errors.New("listing not found")
	// ErrInvalidSellerID = errors.New("invalid seller ID")
	// ErrInvalidDatastreamID = errors.New("invalid datastream ID")
	// ErrListingIssued = errors.New("cannot perform this operation after the listing has been issued")
)

type ContractListingRepository interface {
	BaseRepository[models.ContractListing]
	FindAllBySellerID(sellerID uuid.UUID) ([]models.ContractListing, error)
	FindAllByDatastreamID(datastreamID uuid.UUID) ([]models.ContractListing, error)

	FindAllByMinQuotaReads(minQuota uint64) ([]models.ContractListing, error)
	FindAllByMaxQuotaReads(maxQuota uint64) ([]models.ContractListing, error)
	FindAllByQuotaRange(minQuota, maxQuota uint64) ([]models.ContractListing, error)

	FindAllByMaxReadBytes(maxBytes uint64) ([]models.ContractListing, error)

	FindAllExpiringBefore(deadline time.Time) ([]models.ContractListing, error)
	FindAllValidListings(now time.Time) ([]models.ContractListing, error)
}

type contractListingRepository struct {
	db *gorm.DB
}

func NewContractListingRepository(db *gorm.DB) ContractListingRepository {
	return &contractListingRepository{
		db: db,
	}
}

func manyFromForeign(fk uuid.UUID, db *gorm.DB, reference string) ([]models.ContractListing, error) {
	var listings []models.ContractListing

	result := db.Where(reference+"_id = ?", fk).Find(&listings)
	if result.Error != nil {
		return nil, result.Error
	}
	return listings, nil
}

func NewListingRepository(db *gorm.DB) ContractListingRepository {
	return &contractListingRepository{db: db}
}

func (r *contractListingRepository) FindByID(id uuid.UUID) (*models.ContractListing, error) {
	var listing models.ContractListing
	result := r.db.First(&listing, "listing_id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrListingNotFound
		}
		return nil, result.Error
	}
	return &listing, nil
}

func (r *contractListingRepository) FindAll() ([]models.ContractListing, error) {
	var listings []models.ContractListing

	result := r.db.Find(&listings)
	if result.Error != nil {
		return nil, result.Error
	}
	return listings, nil
}

func (r *contractListingRepository) Create(listing *models.ContractListing) error {
	result := r.db.Create(listing)
	return result.Error
}

func (r *contractListingRepository) Update(listing *models.ContractListing) error {
	result := r.db.Save(listing)
	if result.RowsAffected == 0 {
		return ErrListingNotFound
	}
	return result.Error
}

func (r *contractListingRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.ContractListing{}, "listing_id = ?", id)
	if result.RowsAffected == 0 {
		return ErrListingNotFound
	}
	return result.Error
}

func (r *contractListingRepository) FindAllBySellerID(sellerID uuid.UUID) ([]models.ContractListing, error) {
	return manyFromForeign(sellerID, r.db, "seller")
}

func (r *contractListingRepository) FindAllByDatastreamID(datastreamID uuid.UUID) ([]models.ContractListing, error) {
	return manyFromForeign(datastreamID, r.db, "datastream")
}

func (r *contractListingRepository) FindAllByMinQuotaReads(minQuota uint64) ([]models.ContractListing, error) {
	var listings []models.ContractListing
	result := r.db.Where("quota_reads >= ?", minQuota).Find(&listings)
	if result.Error != nil {
		return nil, result.Error
	}
	return listings, nil
}

func (r *contractListingRepository) FindAllByMaxQuotaReads(maxQuota uint64) ([]models.ContractListing, error) {
	var listings []models.ContractListing
	result := r.db.Where("quota_reads <= ?", maxQuota).Find(&listings)
	if result.Error != nil {
		return nil, result.Error
	}
	return listings, nil
}

func (r *contractListingRepository) FindAllByQuotaRange(minQuota, maxQuota uint64) ([]models.ContractListing, error) {
	var listings []models.ContractListing
	result := r.db.Where("quota_reads BETWEEN ? AND ?", minQuota, maxQuota).Find(&listings)
	if result.Error != nil {
		return nil, result.Error
	}
	return listings, nil
}

func (r *contractListingRepository) FindAllByMaxReadBytes(maxBytes uint64) ([]models.ContractListing, error) {
	var listings []models.ContractListing
	result := r.db.Where("read_bytes <= ?", maxBytes).Find(&listings)
	if result.Error != nil {
		return nil, result.Error
	}
	return listings, nil
}

func (r *contractListingRepository) FindAllExpiringBefore(deadline time.Time) ([]models.ContractListing, error) {
	var listings []models.ContractListing
	result := r.db.Where("exercise_by < ?", deadline).Find(&listings)
	if result.Error != nil {
		return nil, result.Error
	}
	return listings, nil
}

func (r *contractListingRepository) FindAllValidListings(now time.Time) ([]models.ContractListing, error) {
	var listings []models.ContractListing
	result := r.db.Where("exercise_by > ?", now).Find(&listings)
	if result.Error != nil {
		return nil, result.Error
	}
	return listings, nil
}
