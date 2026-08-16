package service

import (
	"errors"
	"net/http"
	"time"

	"keuangan-api/app/model"
	"keuangan-api/app/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService struct {
	Repo *repository.TransactionRepository
}

type CreateTransactionInput struct {
	CategoryID string  `json:"category_id" binding:"required,uuid"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	Date       string  `json:"date" binding:"required"`
	Notes      *string `json:"notes"`
}

type UpdateTransactionInput struct {
	CategoryID string  `json:"category_id" binding:"omitempty,uuid"`
	Amount     float64 `json:"amount" binding:"omitempty,gt=0"`
	Date       string  `json:"date"`
	Notes      *string `json:"notes"`
}

func (s *TransactionService) Create(input CreateTransactionInput, userID uuid.UUID) (*model.Transaction, int, error) {
	categoryID, _ := uuid.Parse(input.CategoryID)
	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("format tanggal tidak valid, gunakan YYYY-MM-DD")
	}

	tx := &model.Transaction{
		UserID:     userID,
		CategoryID: categoryID,
		Amount:     input.Amount,
		Date:       date,
		Notes:      input.Notes,
	}
	if err := s.Repo.Create(tx); err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal menyimpan transaksi")
	}
	return tx, http.StatusCreated, nil
}

func (s *TransactionService) GetByUser(userID uuid.UUID) ([]model.Transaction, int, error) {
	transactions, err := s.Repo.GetByUserID(userID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil transaksi")
	}
	return transactions, http.StatusOK, nil
}

func (s *TransactionService) Update(id uuid.UUID, input UpdateTransactionInput, userID uuid.UUID) (*model.Transaction, int, error) {
	tx, err := s.Repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.New("transaksi tidak ditemukan")
		}
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil transaksi")
	}

	// Hanya pemilik yang boleh edit
	if tx.UserID != userID {
		return nil, http.StatusForbidden, errors.New("kamu tidak memiliki akses ke transaksi ini")
	}

	if input.CategoryID != "" {
		tx.CategoryID, _ = uuid.Parse(input.CategoryID)
	}
	if input.Amount > 0 {
		tx.Amount = input.Amount
	}
	if input.Date != "" {
		date, err := time.Parse("2006-01-02", input.Date)
		if err != nil {
			return nil, http.StatusBadRequest, errors.New("format tanggal tidak valid, gunakan YYYY-MM-DD")
		}
		tx.Date = date
	}
	if input.Notes != nil {
		tx.Notes = input.Notes
	}

	if err := s.Repo.Update(tx); err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengupdate transaksi")
	}
	return tx, http.StatusOK, nil
}

func (s *TransactionService) Delete(id uuid.UUID, userID uuid.UUID) (int, error) {
	tx, err := s.Repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, errors.New("transaksi tidak ditemukan")
		}
		return http.StatusInternalServerError, errors.New("gagal mengambil transaksi")
	}

	if tx.UserID != userID {
		return http.StatusForbidden, errors.New("kamu tidak memiliki akses ke transaksi ini")
	}

	if err := s.Repo.Delete(id); err != nil {
		return http.StatusInternalServerError, errors.New("gagal menghapus transaksi")
	}
	return http.StatusOK, nil
}
