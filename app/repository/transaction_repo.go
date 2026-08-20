package repository

import (
	"time"

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

func (r *TransactionRepository) GetByUserID(userID uuid.UUID, startDate, endDate string, categoryID string, page, limit int) ([]model.Transaction, int64, error) {
	var transactions []model.Transaction
	var total int64

	query := r.DB.Model(&model.Transaction{}).Where("user_id = ?", userID)

	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("date >= ?", t)
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("date <= ?", t)
		}
	}
	if categoryID != "" {
		if catUUID, err := uuid.Parse(categoryID); err == nil {
			query = query.Where("category_id = ?", catUUID)
		}
	}

	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	err := query.Preload("Category").Order("date DESC, created_at DESC").Find(&transactions).Error
	return transactions, total, err
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
	tx.Category = model.Category{} // Clear struct agar tidak ada collision
	if err := r.DB.Omit("Category").Save(tx).Error; err != nil {
		return err
	}
	return r.DB.Preload("Category").First(tx, "id = ?", tx.ID).Error
}

func (r *TransactionRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&model.Transaction{}, "id = ?", id).Error
}

// GetTotalBalance menghitung total saldo keseluruhan (Pemasukan - Pengeluaran) user tanpa batas pagination.
func (r *TransactionRepository) GetTotalBalance(userID uuid.UUID) (float64, error) {
	var totalIncome float64
	var totalExpense float64

	r.DB.Model(&model.Transaction{}).
		Joins("JOIN categories ON categories.id = transactions.category_id").
		Where("transactions.user_id = ? AND categories.type = 'income' AND transactions.deleted_at IS NULL", userID).
		Select("COALESCE(SUM(transactions.amount), 0)").
		Scan(&totalIncome)

	r.DB.Model(&model.Transaction{}).
		Joins("JOIN categories ON categories.id = transactions.category_id").
		Where("transactions.user_id = ? AND categories.type = 'expense' AND transactions.deleted_at IS NULL", userID).
		Select("COALESCE(SUM(transactions.amount), 0)").
		Scan(&totalExpense)

	return totalIncome - totalExpense, nil
}

func (r *TransactionRepository) GetTotalExpenseToday(userID uuid.UUID) (float64, error) {
	var total float64
	todayStr := time.Now().Format("2006-01-02")
	today, _ := time.Parse("2006-01-02", todayStr)

	err := r.DB.Model(&model.Transaction{}).
		Joins("JOIN categories ON categories.id = transactions.category_id").
		Where("transactions.user_id = ? AND transactions.date = ? AND categories.type = 'expense' AND transactions.deleted_at IS NULL", userID, today).
		Select("COALESCE(SUM(transactions.amount), 0)").
		Scan(&total).Error

	return total, err
}

func (r *TransactionRepository) GetTotalExpenseMonth(userID uuid.UUID, month, year int, categoryID *uuid.UUID) (float64, error) {
	var total float64
	query := r.DB.Model(&model.Transaction{}).
		Joins("JOIN categories ON categories.id = transactions.category_id").
		Where("transactions.user_id = ? AND EXTRACT(MONTH FROM transactions.date) = ? AND EXTRACT(YEAR FROM transactions.date) = ? AND categories.type = 'expense' AND transactions.deleted_at IS NULL", userID, month, year)

	if categoryID != nil {
		query = query.Where("transactions.category_id = ?", *categoryID)
	}

	err := query.Select("COALESCE(SUM(transactions.amount), 0)").Scan(&total).Error
	return total, err
}
