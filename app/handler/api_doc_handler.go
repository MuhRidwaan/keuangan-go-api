package handler

import (
	"net/http"
	"os"

	"keuangan-api/app/service"
	"keuangan-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type APIDocHandler struct {
	Service *service.APIDocService
}

// GET /api/docs — Mengambil seluruh data dokumentasi dari DB (JSON)
func (h *APIDocHandler) GetAll(c *gin.Context) {
	docs, code, err := h.Service.GetAll()
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Berhasil mengambil dokumentasi API", docs)
}

// GET /docs — Menampilkan halaman Web UI dokumentasi interaktif
func (h *APIDocHandler) RenderDocsUI(c *gin.Context) {
	filePath := "web/docs.html"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.String(http.StatusNotFound, "File web/docs.html tidak ditemukan")
		return
	}
	c.File(filePath)
}
