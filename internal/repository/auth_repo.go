package repository

import (
	"keuangan-api/internal/model"

	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func (r *AuthRepository) CreateUser(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *AuthRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdatePassword mengupdate password_hash user berdasarkan ID.
func (r *AuthRepository) UpdatePassword(userID string, newHash string) error {
	return r.DB.Model(&model.User{}).
		Where("id = ?", userID).
		Update("password_hash", newHash).Error
}
