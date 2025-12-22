package repos

import (
	"contract_market_demo/backend/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	BaseRepository[models.TransactionRecord]
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) FindByID(id uuid.UUID) (*models.TransactionRecord, error) {
	var record models.TransactionRecord
	result := r.db.First(&record, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrContractHeaderNotFound
		}
		return nil, result.Error
	}
	return &record, nil
}

func (r *transactionRepository) FindAll() ([]models.TransactionRecord, error) {
	var records []models.TransactionRecord
	result := r.db.Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}
	return records, nil
}

func (r *transactionRepository) Create(record *models.TransactionRecord) error {
	result := r.db.Create(record)
	return result.Error
}

func (r *transactionRepository) Update(record *models.TransactionRecord) error {
	result := r.db.Save(record)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrContractHeaderNotFound
	}
	return nil
}

func (r *transactionRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.TransactionRecord{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrContractHeaderNotFound
	}
	return nil
}
