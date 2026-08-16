package handler

import (
	"encoding/json"
	"net/http"
	"sync"

	"keuangan-api/app/database"
	"keuangan-api/app/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	app  *gin.Engine
	once sync.Once
)

func initApp() {
	// 1. Load .env jika ada (lokal), jika tidak ada (Vercel) gunakan env dari OS
	_ = godotenv.Load()

	// 2. Konek ke PostgreSQL
	db := database.Connect()
	if db != nil {
		// 3. Setup router Gin
		app = router.Setup(db)
	}
}

// Handler adalah entrypoint standar Serverless Function untuk Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initApp)

	if app == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"meta": map[string]any{
				"code":    500,
				"status":  "error",
				"message": "Gagal terhubung ke Database PostgreSQL. Pastikan Environment Variables di Vercel Dashboard sudah sesuai.",
			},
			"data": nil,
		})
		return
	}

	app.ServeHTTP(w, r)
}
