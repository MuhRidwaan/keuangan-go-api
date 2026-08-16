package handler

import (
	"net/http"

	"keuangan-api/app/middleware"
	"keuangan-api/app/service"
	"keuangan-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	Service *service.NotificationService
}

// GET /api/notifications
func (h *NotificationHandler) GetAll(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	notifs, code, err := h.Service.GetMyNotifications(userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil notifikasi", notifs)
}

// PUT /api/notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	notifID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID notifikasi tidak valid")
		return
	}

	code, err := h.Service.MarkAsRead(notifID, userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, code, "Notifikasi ditandai sudah dibaca", nil)
}
