package repository

import (
	"keuangan-api/internal/model"

	"gorm.io/gorm"
)

type APIDocRepository struct {
	DB *gorm.DB
}

func (r *APIDocRepository) GetAll() ([]model.APIDocumentation, error) {
	var docs []model.APIDocumentation
	err := r.DB.Order("category ASC, name ASC").Find(&docs).Error
	return docs, err
}
