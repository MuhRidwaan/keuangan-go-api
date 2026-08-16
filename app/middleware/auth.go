package middleware

import (
	"net/http"
	"strings"

	pkgjwt "keuangan-api/pkg/jwt"
	"keuangan-api/pkg/response"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userID"
const UserEmailKey = "userEmail"

// AuthMiddleware memvalidasi JWT dari header Authorization: Bearer <token>.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "Token tidak ditemukan")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := pkgjwt.ValidateToken(tokenStr)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Token tidak valid atau sudah expired")
			c.Abort()
			return
		}

		// Simpan data user ke context agar bisa diakses handler
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)
		c.Next()
	}
}
