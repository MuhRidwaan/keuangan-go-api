package handler

import (
	"strconv"
	"time"

	"keuangan-api/app/service"
	"keuangan-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BudgetHandler struct {
	Service *service.BudgetService
}

func (h *BudgetHandler) Create(c *gin.Context) {
	var input service.CreateBudgetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "Input tidak valid: "+err.Error())
		return
	}

	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	budget, code, err := h.Service.Create(input, userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, code, "Budget bulanan berhasil dibuat", budget)
}

func (h *BudgetHandler) GetMyBudgets(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	monthStr := c.Query("month")
	yearStr := c.Query("year")

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	if month == 0 || year == 0 {
		now := time.Now()
		month = int(now.Month())
		year = now.Year()
	}

	budgets, code, err := h.Service.GetByUserAndPeriod(userID, month, year)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, code, "Berhasil mengambil data budget bulanan", budgets)
}

func (h *BudgetHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(c, 400, "ID budget tidak valid")
		return
	}

	var input service.CreateBudgetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "Input tidak valid: "+err.Error())
		return
	}

	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	budget, code, err := h.Service.Update(id, input, userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, code, "Budget berhasil diperbarui", budget)
}

func (h *BudgetHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(c, 400, "ID budget tidak valid")
		return
	}

	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	code, err := h.Service.Delete(id, userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, code, "Budget berhasil dihapus", nil)
}
