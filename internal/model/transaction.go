package model

import (
	"time"

	"github.com/google/uuid"
)

// Transaction merepresentasikan satu catatan transaksi keuangan milik user.
type Transaction struct {
	BaseModel
	UserID     uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null;index" json:"category_id"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Amount     float64   `gorm:"type:numeric(15,2);not null" json:"amount"`
	Date       time.Time `gorm:"type:date;not null" json:"date"`
	Notes      *string   `gorm:"type:text" json:"notes,omitempty"`
}
