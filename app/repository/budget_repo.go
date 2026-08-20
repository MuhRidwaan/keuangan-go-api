package repository

import (
	"keuangan-api/app/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BudgetRepository struct {
	DB *gorm.DB
}

func (r *BudgetRepository) Create(budget *model.Budget) error {
	if err := r.DB.Create(budget).Error; err != nil {
		return err
	}
	return r.DB.Preload("Category").First(budget, "id = ?", budget.ID).Error
}

func (r *BudgetRepository) GetByUserAndPeriod(userID uuid.UUID, month, year int) ([]model.Budget, error) {
	var budgets []model.Budget
	err := r.DB.Preload("Category").
		Where("user_id = ? AND month = ? AND year = ?", userID, month, year).
		Find(&budgets).Error
	return budgets, err
}

func (r *BudgetRepository) FindByID(id uuid.UUID) (*model.Budget, error) {
	var budget model.Budget
	err := r.DB.Preload("Category").First(&budget, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *BudgetRepository) Update(budget *model.Budget) error {
	budget.Category = nil
	if err := r.DB.Omit("Category").Save(budget).Error; err != nil {
		return err
	}
	return r.DB.Preload("Category").First(budget, "id = ?", budget.ID).Error
}

func (r *BudgetRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&model.Budget{}, "id = ?", id).Error
}
