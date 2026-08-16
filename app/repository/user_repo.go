package repository

import (
	"keuangan-api/app/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.DB.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
