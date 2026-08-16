package handler

import (
	"net/http"

	"keuangan-api/internal/service"
	"keuangan-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	Service *service.AuthService
}

// Register godoc
// POST /api/register
func (h *AuthHandler) Register(c *gin.Context) {
	var input service.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, code, err := h.Service.Register(input)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, code, "Registrasi berhasil", result)
}

// Login godoc
// POST /api/login
func (h *AuthHandler) Login(c *gin.Context) {
	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, code, err := h.Service.Login(input)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, code, "Login berhasil", result)
}

// ChangePassword godoc
// POST /api/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Gagal membaca user ID")
		return
	}

	var input service.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	code, err := h.Service.ChangePassword(userID, input)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, code, "Password berhasil diubah", nil)
}
