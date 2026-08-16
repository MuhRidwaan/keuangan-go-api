package handler

import (
	"net/http"

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

// GET /api/transactions
func (h *TransactionHandler) GetAll(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	transactions, code, err := h.Service.GetByUser(userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Berhasil mengambil transaksi", transactions)
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
