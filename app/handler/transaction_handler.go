package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"keuangan-api/app/middleware"
	"keuangan-api/app/service"
	"keuangan-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	Service *service.TransactionService
}

// POST /api/transactions
func (h *TransactionHandler) Create(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	var input service.CreateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	tx, code, err := h.Service.Create(input, userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Transaksi berhasil disimpan", tx)
}

// GET /api/transactions (Supports filtering & pagination)
func (h *TransactionHandler) GetAll(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	categoryID := c.Query("category_id")

	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	res, code, err := h.Service.GetByUser(userID, startDate, endDate, categoryID, page, limit)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Berhasil mengambil transaksi", res)
}

// GET /api/transactions/report
func (h *TransactionHandler) GetReport(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	res, code, err := h.Service.GetReport(userID, startDate, endDate)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Berhasil mengambil laporan transaksi", res)
}

// GET /api/transactions/export
func (h *TransactionHandler) ExportCSV(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	csvBytes, code, err := h.Service.ExportCSV(userID, startDate, endDate)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	fileName := fmt.Sprintf("laporan_keuangan_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Data(http.StatusOK, "text/csv", csvBytes)
}

// PUT /api/transactions/:id
func (h *TransactionHandler) Update(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID transaksi tidak valid")
		return
	}
	var input service.UpdateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	tx, code, err := h.Service.Update(id, input, userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Transaksi berhasil diupdate", tx)
}

// DELETE /api/transactions/:id
func (h *TransactionHandler) Delete(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID transaksi tidak valid")
		return
	}
	code, err := h.Service.Delete(id, userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Transaksi berhasil dihapus", nil)
}
