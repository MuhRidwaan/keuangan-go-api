package model

import "github.com/google/uuid"

// CategoryType mendefinisikan tipe kategori transaksi.
type CategoryType string

const (
	CategoryIncome  CategoryType = "income"
	CategoryExpense CategoryType = "expense"
)

// Category merepresentasikan kategori transaksi (bisa global atau milik user tertentu).
// UserID nullable berarti kategori bisa bersifat global (system-defined) jika nil.
type Category struct {
	BaseModel
	UserID *uuid.UUID   `gorm:"type:uuid;index" json:"user_id,omitempty"`
	User   *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Name   string       `gorm:"type:varchar(100);not null" json:"name"`
	Type   CategoryType `gorm:"type:varchar(10);not null" json:"type"`
	Icon   *string      `gorm:"type:varchar(100)" json:"icon,omitempty"`
}
