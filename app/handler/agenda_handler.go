package handler

import (
	"net/http"

	"keuangan-api/app/middleware"
	"keuangan-api/app/service"
	"keuangan-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AgendaHandler struct {
	Service *service.AgendaService
}

// POST /api/agendas
func (h *AgendaHandler) CreateAgenda(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	userEmail := c.MustGet(middleware.UserEmailKey).(string)
	var input service.CreateAgendaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	agenda, code, err := h.Service.CreateAgenda(input, userID, userEmail)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Agenda berhasil dibuat", agenda)
}

// GET /api/agendas
func (h *AgendaHandler) GetMyAgendas(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	agendas, code, err := h.Service.GetMyAgendas(userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Berhasil mengambil agenda", agendas)
}

// PUT /api/agendas/:id
func (h *AgendaHandler) UpdateAgenda(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	agendaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID agenda tidak valid")
		return
	}
	var input service.UpdateAgendaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	agenda, code, err := h.Service.Update(agendaID, input, userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Agenda berhasil diupdate", agenda)
}

// PUT /api/agendas/:id/status
func (h *AgendaHandler) UpdateStatus(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	agendaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID agenda tidak valid")
		return
	}
	var input service.UpdateAgendaStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	agenda, code, err := h.Service.UpdateStatus(agendaID, input, userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Status agenda berhasil diperbarui", agenda)
}

// DELETE /api/agendas/:id
func (h *AgendaHandler) DeleteAgenda(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	agendaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID agenda tidak valid")
		return
	}
	code, err := h.Service.Delete(agendaID, userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Agenda berhasil dihapus", nil)
}

// POST /api/agendas/:id/members
func (h *AgendaHandler) AddMember(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	userEmail := c.MustGet(middleware.UserEmailKey).(string)
	agendaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID agenda tidak valid")
		return
	}
	var input service.AddAgendaMemberInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	member, code, err := h.Service.AddMember(agendaID, input, userID, userEmail)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Anggota berhasil ditambahkan ke agenda", member)
}
