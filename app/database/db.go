package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect membuka koneksi ke PostgreSQL menggunakan GORM dan menjalankan AutoMigrate.
func Connect() *gorm.DB {
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Tampilkan query SQL di log (ubah ke Silent di production)
	})
	if err != nil {
		log.Fatalf("Gagal konek ke database: %v", err)
	}

	fmt.Println("Berhasil terhubung ke PostgreSQL via GORM!")
	return db
}
