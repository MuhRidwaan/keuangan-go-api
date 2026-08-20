package handler

import (
	"net/http"

	"keuangan-api/app/middleware"
	"keuangan-api/app/service"
	"keuangan-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	Service *service.UserService
}

// GET /api/contacts
func (h *UserHandler) GetRecentContacts(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	contacts, code, err := h.Service.GetRecentContacts(userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Berhasil mengambil kontak riwayat", contacts)
}
