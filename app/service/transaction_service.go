package service

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"time"

	"keuangan-api/app/model"
	"keuangan-api/app/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService struct {
	Repo       *repository.TransactionRepository
	BudgetRepo *repository.BudgetRepository
	NotifRepo  *repository.NotificationRepository
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

type TransactionListResponse struct {
	Items             []model.Transaction `json:"items"`
	TotalItems        int64               `json:"total_items"`
	Page              int                 `json:"page"`
	Limit             int                 `json:"limit"`
	TotalPages        int                 `json:"total_pages"`
	TotalExpenseToday float64             `json:"total_expense_today"`
}

type CategoryReport struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	CategoryType string  `json:"category_type"`
	TotalAmount  float64 `json:"total_amount"`
	Percentage   float64 `json:"percentage"`
}

type ReportSummaryResponse struct {
	StartDate     string           `json:"start_date"`
	EndDate       string           `json:"end_date"`
	TotalIncome   float64          `json:"total_income"`
	TotalExpense  float64          `json:"total_expense"`
	NetBalance    float64          `json:"net_balance"`
	Breakdown     []CategoryReport `json:"breakdown"`
	TransactionCount int           `json:"transaction_count"`
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

	// Trigger pengecekan budget limit secara otomatis
	s.checkBudgetLimit(userID, date, categoryID)

	return tx, http.StatusCreated, nil
}

func (s *TransactionService) GetByUser(userID uuid.UUID, startDate, endDate, categoryID string, page, limit int) (*TransactionListResponse, int, error) {
	transactions, totalItems, err := s.Repo.GetByUserID(userID, startDate, endDate, categoryID, page, limit)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil transaksi")
	}

	totalExpenseToday, _ := s.Repo.GetTotalExpenseToday(userID)

	totalPages := 1
	if limit > 0 && totalItems > 0 {
		totalPages = int((totalItems + int64(limit) - 1) / int64(limit))
	}

	res := &TransactionListResponse{
		Items:             transactions,
		TotalItems:        totalItems,
		Page:              page,
		Limit:             limit,
		TotalPages:        totalPages,
		TotalExpenseToday: totalExpenseToday,
	}

	return res, http.StatusOK, nil
}

func (s *TransactionService) Update(id uuid.UUID, input UpdateTransactionInput, userID uuid.UUID) (*model.Transaction, int, error) {
	tx, err := s.Repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.New("transaksi tidak ditemukan")
		}
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil transaksi")
	}

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

	// Trigger pengecekan budget limit
	s.checkBudgetLimit(userID, tx.Date, tx.CategoryID)

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

func (s *TransactionService) GetReport(userID uuid.UUID, startDate, endDate string) (*ReportSummaryResponse, int, error) {
	transactions, _, err := s.Repo.GetByUserID(userID, startDate, endDate, "", 0, 0)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil data laporan")
	}

	var totalIncome, totalExpense float64
	categoryTotals := make(map[string]*CategoryReport)

	for _, tx := range transactions {
		catName := tx.Category.Name
		catType := string(tx.Category.Type)
		if catName == "" {
			catName = "Tanpa Kategori"
			catType = "expense"
		}

		if catType == "income" {
			totalIncome += tx.Amount
		} else {
			totalExpense += tx.Amount
		}

		if _, exists := categoryTotals[catName]; !exists {
			categoryTotals[catName] = &CategoryReport{
				CategoryID:   tx.CategoryID.String(),
				CategoryName: catName,
				CategoryType: catType,
				TotalAmount:  0,
			}
		}
		categoryTotals[catName].TotalAmount += tx.Amount
	}

	breakdown := make([]CategoryReport, 0, len(categoryTotals))
	for _, item := range categoryTotals {
		if item.CategoryType == "expense" && totalExpense > 0 {
			item.Percentage = (item.TotalAmount / totalExpense) * 100
		} else if item.CategoryType == "income" && totalIncome > 0 {
			item.Percentage = (item.TotalAmount / totalIncome) * 100
		}
		breakdown = append(breakdown, *item)
	}

	return &ReportSummaryResponse{
		StartDate:        startDate,
		EndDate:          endDate,
		TotalIncome:      totalIncome,
		TotalExpense:     totalExpense,
		NetBalance:       totalIncome - totalExpense,
		Breakdown:        breakdown,
		TransactionCount: len(transactions),
	}, http.StatusOK, nil
}

func (s *TransactionService) ExportCSV(userID uuid.UUID, startDate, endDate string) ([]byte, int, error) {
	transactions, _, err := s.Repo.GetByUserID(userID, startDate, endDate, "", 0, 0)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil data untuk ekspor")
	}

	var b bytes.Buffer
	writer := csv.NewWriter(&b)

	// Write Header
	_ = writer.Write([]string{"ID", "Tanggal", "Kategori", "Tipe", "Jumlah (Rp)", "Catatan"})

	for _, tx := range transactions {
		catName := tx.Category.Name
		catType := string(tx.Category.Type)
		notes := ""
		if tx.Notes != nil {
			notes = *tx.Notes
		}

		_ = writer.Write([]string{
			tx.ID.String(),
			tx.Date.Format("2006-01-02"),
			catName,
			catType,
			fmt.Sprintf("%.2f", tx.Amount),
			notes,
		})
	}

	writer.Flush()
	return b.Bytes(), http.StatusOK, nil
}

func (s *TransactionService) checkBudgetLimit(userID uuid.UUID, date time.Time, categoryID uuid.UUID) {
	if s.BudgetRepo == nil || s.NotifRepo == nil {
		return
	}

	month := int(date.Month())
	year := date.Year()

	budgets, err := s.BudgetRepo.GetByUserAndPeriod(userID, month, year)
	if err != nil || len(budgets) == 0 {
		return
	}

	for _, b := range budgets {
		// Jika budget khusus kategori dan tidak cocok, skip
		if b.CategoryID != nil && *b.CategoryID != categoryID {
			continue
		}

		totalExpense, err := s.Repo.GetTotalExpenseMonth(userID, month, year, b.CategoryID)
		if err != nil || b.Amount <= 0 {
			continue
		}

		pct := (totalExpense / b.Amount) * 100.0

		if pct >= 100.0 {
			_ = s.NotifRepo.Create(&model.Notification{
				UserID:  userID,
				Title:   "⚠️ Over Budget!",
				Message: fmt.Sprintf("Pengeluaran Anda bulan ini (Rp %.0f) telah MELEBIHI batas budget sebesar Rp %.0f!", totalExpense, b.Amount),
			})
		} else if pct >= 80.0 {
			_ = s.NotifRepo.Create(&model.Notification{
				UserID:  userID,
				Title:   "🔔 Peringatan Budget!",
				Message: fmt.Sprintf("Pengeluaran Anda bulan ini (Rp %.0f) sudah mencapai %.1f%% dari batas budget Rp %.0f.", totalExpense, pct, b.Amount),
			})
		}
	}
}
