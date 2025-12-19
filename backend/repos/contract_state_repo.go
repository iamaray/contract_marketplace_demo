package repos

import (
	"errors"
	"contract_market_demo/backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrContractStateNotFound = errors.New("contract state not found")
)

type ContractStateRepository interface {
	BaseRepository[models.ContractState]
	FindAllByOwnerID(ownerID uuid.UUID) ([]models.ContractState, error)

	FindAllByStatus(status models.ContractStatus) ([]models.ContractState, error)
	FindByContractIDAndStatus(contractID uuid.UUID, status models.ContractStatus) (*models.ContractState, error)

	FindByContractID(contractID uuid.UUID) (*models.ContractState, error)

	UpdateStatus(contractID uuid.UUID, newStatus models.ContractStatus) error
	UpdateReadsRemaining(contractID uuid.UUID, readsRemaining uint64) error
}

type contractStateRepository struct {
	db *gorm.DB
}

func NewContractStateRepository(db *gorm.DB) ContractStateRepository {
	return &contractStateRepository{db: db}
}

func (r *contractStateRepository) FindByID(id uuid.UUID) (*models.ContractState, error) {
	var contractState models.ContractState
	result := r.db.First(&contractState, "header_id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrContractStateNotFound
		}
		return nil, result.Error
	}
	return &contractState, nil
}

func (r *contractStateRepository) FindAll() ([]models.ContractState, error) {
	var contractStates []models.ContractState
	result := r.db.Find(&contractStates)
	if result.Error != nil {
		return nil, result.Error
	}
	return contractStates, nil
}

func (r *contractStateRepository) Create(contractState *models.ContractState) error {
	result := r.db.Create(contractState)
	return result.Error
}

func (r *contractStateRepository) Update(contractState *models.ContractState) error {
	result := r.db.Save(contractState)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrContractStateNotFound
	}
	return nil
}

func (r *contractStateRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.ContractState{}, "contract_id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrContractStateNotFound
	}
	return nil
}

func (r *contractStateRepository) FindAllByOwnerID(ownerID uuid.UUID) ([]models.ContractState, error) {
	var contractStates []models.ContractState
	result := r.db.Where("owner_id = ?", ownerID).Find(&contractStates)
	if result.Error != nil {
		return nil, result.Error
	}
	return contractStates, nil
}

func (r *contractStateRepository) FindAllByStatus(status models.ContractStatus) ([]models.ContractState, error) {
	var contractStates []models.ContractState
	result := r.db.Where("status = ?", status).Find(&contractStates)
	if result.Error != nil {
		return nil, result.Error
	}
	return contractStates, nil
}

func (r *contractStateRepository) FindByContractID(contractID uuid.UUID) (*models.ContractState, error) {
	var contractState models.ContractState
	result := r.db.Where("contract_id = ?", contractID).First(&contractState)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrContractStateNotFound
		}
		return nil, result.Error
	}
	return &contractState, nil
}

func (r *contractStateRepository) FindByContractIDAndStatus(contractID uuid.UUID, status models.ContractStatus) (*models.ContractState, error) {
	var contractState models.ContractState
	result := r.db.Where("contract_id = ? AND status = ?", contractID, status).First(&contractState)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrContractStateNotFound
		}
		return nil, result.Error
	}
	return &contractState, nil
}

func (r *contractStateRepository) UpdateStatus(contractID uuid.UUID, newStatus models.ContractStatus) error {
	result := r.db.Model(&models.ContractState{}).
		Where("contract_id = ?", contractID).
		Update("status", newStatus)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrContractStateNotFound
	}
	return nil
}

func (r *contractStateRepository) UpdateReadsRemaining(contractID uuid.UUID, readsRemaining uint64) error {
	result := r.db.Model(&models.ContractState{}).
		Where("contract_id = ?", contractID).
		Update("reads_remaining", readsRemaining)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrContractStateNotFound
	}
	return nil
}
