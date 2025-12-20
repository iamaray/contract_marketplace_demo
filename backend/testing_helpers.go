package main

import (
	"contract_market_demo/backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateDummySellers(n int, db *gorm.DB) ([]*models.User, error) {
	users := make([]*models.User, n)
	for i := 0; i < n; i++ {
		new_user := &models.User{ID: uuid.New()}
		users[i] = new_user

		result := db.Create(new_user)
		err := result.Error
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}
