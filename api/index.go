package handler

import (
	"net/http"
	"sync"

	"keuangan-api/internal/database"
	"keuangan-api/internal/router"

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

	// 3. Setup router Gin
	app = router.Setup(db)
}

// Handler adalah entrypoint standar Serverless Function untuk Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initApp)
	app.ServeHTTP(w, r)
}
