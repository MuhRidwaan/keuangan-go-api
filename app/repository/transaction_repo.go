package repository

import (
	"keuangan-api/app/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	DB *gorm.DB
}

func (r *TransactionRepository) Create(tx *model.Transaction) error {
	if err := r.DB.Create(tx).Error; err != nil {
		return err
	}
	return r.DB.Preload("Category").First(tx, "id = ?", tx.ID).Error
}

func (r *TransactionRepository) GetByUserID(userID uuid.UUID) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.DB.
		Preload("Category").
		Where("user_id = ?", userID).
		Order("date DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) FindByID(id uuid.UUID) (*model.Transaction, error) {
	var tx model.Transaction
	err := r.DB.Preload("Category").First(&tx, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *TransactionRepository) Update(tx *model.Transaction) error {
	if err := r.DB.Save(tx).Error; err != nil {
		return err
	}
	return r.DB.Preload("Category").First(tx, "id = ?", tx.ID).Error
}

// Delete melakukan soft delete (set deleted_at).
func (r *TransactionRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&model.Transaction{}, "id = ?", id).Error
}
