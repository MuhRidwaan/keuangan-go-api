package model

import "github.com/google/uuid"

// Budget merepresentasikan batas anggaran bulanan pengguna.
type Budget struct {
	BaseModel
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	User       User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CategoryID *uuid.UUID `gorm:"type:uuid;index" json:"category_id,omitempty"`
	Category   *Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Amount     float64    `gorm:"type:numeric(15,2);not null" json:"amount"`
	Month      int        `gorm:"not null" json:"month"`
	Year       int        `gorm:"not null" json:"year"`
}
