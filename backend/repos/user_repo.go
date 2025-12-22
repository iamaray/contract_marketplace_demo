package repos

import (
	"contract_market_demo/backend/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	BaseRepository[models.User]
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrContractHeaderNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *userRepository) Create(user *models.User) error {
	result := r.db.Create(user)
	return result.Error
}

func (r *userRepository) Update(user *models.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrContractHeaderNotFound
	}
	return nil
}

func (r *userRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrContractHeaderNotFound
	}
	return nil
}
