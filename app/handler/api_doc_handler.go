package handler

import (
	_ "embed"
	"net/http"

	"keuangan-api/app/service"
	"keuangan-api/pkg/response"

	"github.com/gin-gonic/gin"
)

//go:embed docs.html
var docsHTML string

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
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, docsHTML)
}
