package handler

import (
	"net/http"

	"keuangan-api/internal/middleware"
	"keuangan-api/internal/service"
	"keuangan-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SavingHandler struct {
	Service *service.SavingService
}

// POST /api/savings
func (h *SavingHandler) CreateGoal(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	userEmail := c.MustGet(middleware.UserEmailKey).(string)
	var input service.CreateGoalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	goal, code, err := h.Service.CreateGoal(input, userID, userEmail)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Saving goal berhasil dibuat", goal)
}

// GET /api/savings
func (h *SavingHandler) GetMyGoals(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	goals, code, err := h.Service.GetMyGoals(userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Berhasil mengambil saving goals", goals)
}

// PUT /api/savings/:id
func (h *SavingHandler) UpdateGoal(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID goal tidak valid")
		return
	}
	var input service.UpdateGoalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	goal, code, err := h.Service.Update(goalID, input, userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Saving goal berhasil diupdate", goal)
}

// DELETE /api/savings/:id
func (h *SavingHandler) DeleteGoal(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID goal tidak valid")
		return
	}
	code, err := h.Service.Delete(goalID, userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Saving goal berhasil dihapus", nil)
}

// POST /api/savings/:id/members
func (h *SavingHandler) AddMember(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	userEmail := c.MustGet(middleware.UserEmailKey).(string)
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID goal tidak valid")
		return
	}
	var input service.AddMemberByEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	member, code, err := h.Service.AddMember(goalID, input, userID, userEmail)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Anggota berhasil ditambahkan", member)
}

// POST /api/savings/:id/contribute
func (h *SavingHandler) Contribute(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID goal tidak valid")
		return
	}
	var input service.ContributeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	contribution, code, err := h.Service.Contribute(goalID, input, userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Kontribusi berhasil disimpan", contribution)
}

// POST /api/savings/:id/withdraw
func (h *SavingHandler) Withdraw(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID goal tidak valid")
		return
	}
	var input service.WithdrawInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	contribution, code, err := h.Service.Withdraw(goalID, input, userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Penarikan berhasil diproses", contribution)
}

// GET /api/savings/:id/contributions
func (h *SavingHandler) GetContributionHistory(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID goal tidak valid")
		return
	}

	result, code, err := h.Service.GetContributionHistory(goalID, userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil riwayat kontribusi", result)
}
