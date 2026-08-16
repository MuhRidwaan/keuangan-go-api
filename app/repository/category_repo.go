package repository

import (
	"keuangan-api/app/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

// GetVisibleCategories mengambil kategori sistem (user_id NULL) + kategori milik user.
func (r *CategoryRepository) GetVisibleCategories(userID uuid.UUID) ([]model.Category, error) {
	var categories []model.Category
	err := r.DB.Where("user_id IS NULL OR user_id = ?", userID).Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) Create(category *model.Category) error {
	return r.DB.Create(category).Error
}

func (r *CategoryRepository) FindByID(id uuid.UUID) (*model.Category, error) {
	var category model.Category
	err := r.DB.First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}
