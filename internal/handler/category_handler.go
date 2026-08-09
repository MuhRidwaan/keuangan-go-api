package handler

import (
	"net/http"

	"keuangan-api/internal/middleware"
	"keuangan-api/internal/service"
	"keuangan-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	Service *service.CategoryService
}

// GetCategories godoc
// GET /api/categories
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	categories, code, err := h.Service.GetCategories(userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil kategori", categories)
}

// CreateCategory godoc
// POST /api/categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	var input service.CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	category, code, err := h.Service.CreateCategory(input, userID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, code, "Kategori berhasil dibuat", category)
}
