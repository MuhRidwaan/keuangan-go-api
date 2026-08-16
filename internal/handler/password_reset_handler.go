package handler

import (
	"net/http"

	"keuangan-api/internal/service"
	"keuangan-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type PasswordResetHandler struct {
	Service *service.PasswordResetService
}

// POST /api/forgot-password
func (h *PasswordResetHandler) ForgotPassword(c *gin.Context) {
	var input service.ForgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	token, code, err := h.Service.ForgotPassword(input)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Jika email terdaftar, token reset password telah dibuat", map[string]string{"token": token})
}

// POST /api/reset-password
func (h *PasswordResetHandler) ResetPassword(c *gin.Context) {
	var input service.ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	code, err := h.Service.ResetPassword(input)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Password berhasil direset, silakan login kembali", nil)
}
