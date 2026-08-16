package main

import (
	"fmt"
	"log"
	"os"

	"keuangan-api/app/database"
	"keuangan-api/app/router"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Load environment variables dari .env jika ada (diabaikan jika di Vercel/production)
	_ = godotenv.Load()

	// 2. Konek ke PostgreSQL via GORM
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Gagal konek ke database: %v", err)
	}

	// Pastikan koneksi ditutup saat aplikasi mati
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// 3. Setup semua routes
	r := router.Setup(db)

	// 4. Jalankan server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server berjalan di http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
