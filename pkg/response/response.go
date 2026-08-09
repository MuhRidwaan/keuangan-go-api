package response

import "github.com/gin-gonic/gin"

// Meta berisi metadata standar untuk setiap response API.
type Meta struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// APIResponse adalah struktur response standar seluruh API.
type APIResponse struct {
	Meta Meta `json:"meta"`
	Data any  `json:"data"`
}

func Success(c *gin.Context, httpCode int, message string, data any) {
	c.JSON(httpCode, APIResponse{
		Meta: Meta{Code: httpCode, Status: "success", Message: message},
		Data: data,
	})
}

func Error(c *gin.Context, httpCode int, message string) {
	c.JSON(httpCode, APIResponse{
		Meta: Meta{Code: httpCode, Status: "error", Message: message},
		Data: nil,
	})
}
