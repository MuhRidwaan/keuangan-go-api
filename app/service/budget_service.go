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

type BudgetService struct {
	Repo            *repository.BudgetRepository
	TransactionRepo *repository.TransactionRepository
}

type CreateBudgetInput struct {
	CategoryID *string `json:"category_id"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	Month      int     `json:"month" binding:"required,min=1,max=12"`
	Year       int     `json:"year" binding:"required,min=2020"`
}

type BudgetWithUsage struct {
	model.Budget
	UsedAmount      float64 `json:"used_amount"`
	RemainingAmount float64 `json:"remaining_amount"`
	Percentage      float64 `json:"percentage"`
}

func (s *BudgetService) Create(input CreateBudgetInput, userID uuid.UUID) (*model.Budget, int, error) {
	var catUUID *uuid.UUID
	if input.CategoryID != nil && *input.CategoryID != "" {
		parsed, err := uuid.Parse(*input.CategoryID)
		if err != nil {
			return nil, http.StatusBadRequest, errors.New("category_id tidak valid")
		}
		catUUID = &parsed
	}

	budget := &model.Budget{
		UserID:     userID,
		CategoryID: catUUID,
		Amount:     input.Amount,
		Month:      input.Month,
		Year:       input.Year,
	}

	if err := s.Repo.Create(budget); err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal menyimpan budget")
	}
	return budget, http.StatusCreated, nil
}

func (s *BudgetService) GetByUserAndPeriod(userID uuid.UUID, month, year int) ([]BudgetWithUsage, int, error) {
	if month == 0 || year == 0 {
		now := time.Now()
		month = int(now.Month())
		year = now.Year()
	}

	budgets, err := s.Repo.GetByUserAndPeriod(userID, month, year)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil data budget")
	}

	result := make([]BudgetWithUsage, 0, len(budgets))
	for _, b := range budgets {
		used, _ := s.TransactionRepo.GetTotalExpenseMonth(userID, b.Month, b.Year, b.CategoryID)
		remaining := b.Amount - used
		if remaining < 0 {
			remaining = 0
		}
		percentage := 0.0
		if b.Amount > 0 {
			percentage = (used / b.Amount) * 100.0
		}

		result = append(result, BudgetWithUsage{
			Budget:          b,
			UsedAmount:      used,
			RemainingAmount: remaining,
			Percentage:      percentage,
		})
	}

	return result, http.StatusOK, nil
}

func (s *BudgetService) Update(id uuid.UUID, input CreateBudgetInput, userID uuid.UUID) (*model.Budget, int, error) {
	budget, err := s.Repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.New("budget tidak ditemukan")
		}
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil data budget")
	}

	if budget.UserID != userID {
		return nil, http.StatusForbidden, errors.New("kamu tidak memiliki akses ke budget ini")
	}

	if input.CategoryID != nil && *input.CategoryID != "" {
		parsed, _ := uuid.Parse(*input.CategoryID)
		budget.CategoryID = &parsed
	} else {
		budget.CategoryID = nil
	}
	budget.Amount = input.Amount
	budget.Month = input.Month
	budget.Year = input.Year

	if err := s.Repo.Update(budget); err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengupdate budget")
	}

	return budget, http.StatusOK, nil
}

func (s *BudgetService) Delete(id uuid.UUID, userID uuid.UUID) (int, error) {
	budget, err := s.Repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, errors.New("budget tidak ditemukan")
		}
		return http.StatusInternalServerError, errors.New("gagal mengambil data budget")
	}

	if budget.UserID != userID {
		return http.StatusForbidden, errors.New("kamu tidak memiliki akses ke budget ini")
	}

	if err := s.Repo.Delete(id); err != nil {
		return http.StatusInternalServerError, errors.New("gagal menghapus budget")
	}
	return http.StatusOK, nil
}
