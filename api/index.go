package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"keuangan-api/app/database"
	"keuangan-api/app/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	app     *gin.Engine
	once    sync.Once
	dbErr   error
)

func initApp() {
	// 1. Load .env jika ada (lokal)
	_ = godotenv.Load()

	// 2. Konek ke PostgreSQL
	db, err := database.Connect()
	if err != nil {
		dbErr = err
		return
	}

	// 3. Setup router Gin
	app = router.Setup(db)
}

// Handler adalah entrypoint standar Serverless Function untuk Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initApp)

	if app == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		errMsg := "Gagal terhubung ke Database PostgreSQL."
		if dbErr != nil {
			errMsg = fmt.Sprintf("Gagal konek ke DB [%s:%s]: %v", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), dbErr)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"meta": map[string]any{
				"code":    500,
				"status":  "error",
				"message": errMsg,
			},
			"data": nil,
		})
		return
	}

	app.ServeHTTP(w, r)
}
