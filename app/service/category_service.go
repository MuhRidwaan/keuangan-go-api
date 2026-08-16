package service

import (
	"errors"
	"net/http"

	"keuangan-api/app/model"
	"keuangan-api/app/repository"

	"github.com/google/uuid"
)

type CategoryService struct {
	Repo *repository.CategoryRepository
}

// CreateCategoryInput adalah DTO untuk membuat kategori baru.
type CreateCategoryInput struct {
	Name string           `json:"name" binding:"required"`
	Type model.CategoryType `json:"type" binding:"required,oneof=income expense"`
	Icon *string          `json:"icon"`
}

func (s *CategoryService) GetCategories(userID uuid.UUID) ([]model.Category, int, error) {
	categories, err := s.Repo.GetVisibleCategories(userID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil kategori")
	}
	return categories, http.StatusOK, nil
}

func (s *CategoryService) CreateCategory(input CreateCategoryInput, userID uuid.UUID) (*model.Category, int, error) {
	category := &model.Category{
		UserID: &userID, // selalu set ke user yang sedang login
		Name:   input.Name,
		Type:   input.Type,
		Icon:   input.Icon,
	}

	if err := s.Repo.Create(category); err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal menyimpan kategori")
	}

	return category, http.StatusCreated, nil
}
