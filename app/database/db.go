package database

import (
	"fmt"
	"log"
	"os"

	"keuangan-api/app/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect membuka koneksi ke PostgreSQL menggunakan GORM.
func Connect() (*gorm.DB, error) {
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

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Disables prepared statement caching collisions
	}), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: false,
	})
	if err != nil {
		log.Printf("ERROR: Gagal konek ke database: %v", err)
		return nil, err
	}

	// Ensure missing columns like saving_contributions.type are patched automatically
	_ = db.Exec(`ALTER TABLE saving_contributions ADD COLUMN IF NOT EXISTS type VARCHAR(3) NOT NULL DEFAULT 'in';`).Error

	// AutoMigrate model agar kolom baru ter-sync ke database
	if err := db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Transaction{},
		&model.SavingGoal{},
		&model.SavingMember{},
		&model.SavingContribution{},
		&model.Agenda{},
		&model.AgendaMember{},
		&model.Notification{},
	); err != nil {
		log.Printf("Warning AutoMigrate: %v", err)
	}

	fmt.Println("Berhasil terhubung ke PostgreSQL via GORM!")
	return db, nil
}
