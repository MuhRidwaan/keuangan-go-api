package main

import (
	"fmt"
	"log"
	"os"

	"keuangan-api/internal/database"
	"keuangan-api/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Load environment variables dari .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Konek ke PostgreSQL via GORM + jalankan AutoMigrate
	db := database.Connect()

	// Pastikan koneksi ditutup saat aplikasi mati
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// 3. Setup semua routes (DI dilakukan di dalam router.Setup)
	r := router.Setup(db)

	// 4. Jalankan server
	port := os.Getenv("PORT")
	fmt.Printf("Server berjalan di http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
