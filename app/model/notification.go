package model

import (
	"time"

	"github.com/google/uuid"
)

// Notification merepresentasikan satu notifikasi untuk seorang user.
type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	IsRead    bool      `gorm:"not null;default:false" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
