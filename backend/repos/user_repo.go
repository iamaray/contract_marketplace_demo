package repos

import (
	"contract_market_demo/backend/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("User not found")

type UserRepository interface {
	BaseRepository[models.User]
	FindByAuth(provider, subject string) (*models.User, error)
	FindOrCreateByAuth(provider, subject, email string) (*models.User, error)
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

func (r *userRepository) FindByAuth(provider, subject string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "auth_provider = ? AND auth_subject = ?", provider, subject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) FindOrCreateByAuth(provider, subject, email string) (*models.User, error) {
	u, err := r.FindByAuth(provider, subject)
	if err == nil {
		if u.Email == "" && email != "" {
			u.Email = email
			_ = r.Update(u)
		}
		return u, nil
	}
	if !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		AuthProvider: provider,
		AuthSubject:  subject,
	}
	if err := r.Create(user); err != nil {
		return nil, err
	}
	return user, nil
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
