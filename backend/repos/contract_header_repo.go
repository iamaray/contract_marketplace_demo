package repos

import (
	"contract_market_demo/backend/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrContractHeaderNotFound = errors.New("contract header not found")
)

type ContractHeaderRepository interface {
	BaseRepository[models.ContractHeader]
	FindAllBySellerID(sellerID uuid.UUID) ([]models.ContractHeader, error)
	FindAllByDatastreamID(datastreamID uuid.UUID) ([]models.ContractHeader, error)
	FindAllByListingID(listingID uuid.UUID) ([]models.ContractHeader, error)

	FindAllByRemainingQuota(minRemainingQuota uint64) ([]models.ContractHeader, error)
	FindAllByQuotaRange(minQuota, maxQuota uint64) ([]models.ContractHeader, error)

	FindAllByReadBytes(maxBytes uint64) ([]models.ContractHeader, error)

	FindAllExpiringBefore(deadline time.Time) ([]models.ContractHeader, error)
	FindAllExpired(now time.Time) ([]models.ContractHeader, error)
	FindAllActive(now time.Time) ([]models.ContractHeader, error)

	FindAllByPriceRange(minPriceNanos, maxPriceNanos int64) ([]models.ContractHeader, error)
}

type contractHeaderRepository struct {
	db *gorm.DB
}

func NewContractHeaderRepository(db *gorm.DB) ContractHeaderRepository {
	return &contractHeaderRepository{db: db}
}

func (r *contractHeaderRepository) FindByID(id uuid.UUID) (*models.ContractHeader, error) {
	var contractHeader models.ContractHeader
	result := r.db.First(&contractHeader, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrContractHeaderNotFound
		}
		return nil, result.Error
	}
	return &contractHeader, nil
}

func (r *contractHeaderRepository) FindAll() ([]models.ContractHeader, error) {
	var contractHeaders []models.ContractHeader
	result := r.db.Find(&contractHeaders)
	if result.Error != nil {
		return nil, result.Error
	}
	return contractHeaders, nil
}

func (r *contractHeaderRepository) Create(contractHeader *models.ContractHeader) error {
	result := r.db.Create(contractHeader)
	return result.Error
}

func (r *contractHeaderRepository) Update(contractHeader *models.ContractHeader) error {
	result := r.db.Save(contractHeader)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrContractHeaderNotFound
	}
	return nil
}

func (r *contractHeaderRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.ContractHeader{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrContractHeaderNotFound
	}
	return nil
}

func (r *contractHeaderRepository) FindAllBySellerID(sellerID uuid.UUID) ([]models.ContractHeader, error) {
	var contractHeaders []models.ContractHeader
	result := r.db.Where("seller_id = ?", sellerID).Find(&contractHeaders)
	if result.Error != nil {
		return nil, result.Error
	}
	return contractHeaders, nil
}

func (r *contractHeaderRepository) FindAllByDatastreamID(datastreamID uuid.UUID) ([]models.ContractHeader, error) {
	var contractHeaders []models.ContractHeader
	result := r.db.Where("datastream_id = ?", datastreamID).Find(&contractHeaders)
	if result.Error != nil {
		return nil, result.Error
	}
	return contractHeaders, nil
}

func (r *contractHeaderRepository) FindAllByListingID(listingID uuid.UUID) ([]models.ContractHeader, error) {
	var contractHeaders []models.ContractHeader
	result := r.db.Where("listing_id = ?", listingID).Find(&contractHeaders)
	if result.Error != nil {
		return nil, result.Error
	}
	return contractHeaders, nil
}

func (r *contractHeaderRepository) FindAllByRemainingQuota(minRemainingQuota uint64) ([]models.ContractHeader, error) {
	var headers []models.ContractHeader
	result := r.db.Where("quota_reads >= ?", minRemainingQuota).Find(&headers)
	if result.Error != nil {
		return nil, result.Error
	}
	return headers, nil
}

func (r *contractHeaderRepository) FindAllByQuotaRange(minQuota, maxQuota uint64) ([]models.ContractHeader, error) {
	var headers []models.ContractHeader
	result := r.db.Where("quota_reads BETWEEN ? AND ?", minQuota, maxQuota).Find(&headers)
	if result.Error != nil {
		return nil, result.Error
	}
	return headers, nil
}

func (r *contractHeaderRepository) FindAllByReadBytes(maxBytes uint64) ([]models.ContractHeader, error) {
	var headers []models.ContractHeader
	result := r.db.Where("read_bytes <= ?", maxBytes).Find(&headers)
	if result.Error != nil {
		return nil, result.Error
	}
	return headers, nil
}

func (r *contractHeaderRepository) FindAllExpiringBefore(deadline time.Time) ([]models.ContractHeader, error) {
	var headers []models.ContractHeader
	result := r.db.Where("exercise_by < ?", deadline).Find(&headers)
	if result.Error != nil {
		return nil, result.Error
	}
	return headers, nil
}

func (r *contractHeaderRepository) FindAllExpired(now time.Time) ([]models.ContractHeader, error) {
	var headers []models.ContractHeader
	result := r.db.Where("exercise_by <= ?", now).Find(&headers)
	if result.Error != nil {
		return nil, result.Error
	}
	return headers, nil
}

func (r *contractHeaderRepository) FindAllActive(now time.Time) ([]models.ContractHeader, error) {
	var headers []models.ContractHeader
	result := r.db.Where("exercise_by > ?", now).Find(&headers)
	if result.Error != nil {
		return nil, result.Error
	}
	return headers, nil
}

func (r *contractHeaderRepository) FindAllByPriceRange(minPriceNanos, maxPriceNanos int64) ([]models.ContractHeader, error) {
	var headers []models.ContractHeader
	result := r.db.Where("list_price_nanos BETWEEN ? AND ?", minPriceNanos, maxPriceNanos).Find(&headers)
	if result.Error != nil {
		return nil, result.Error
	}
	return headers, nil
}
