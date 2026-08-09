package model

import (
	"time"

	"github.com/google/uuid"
)

// PasswordResetToken menyimpan token reset password yang dikirim ke email user.
type PasswordResetToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	UsedAt    *time.Time `gorm:"default:null" json:"used_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// IsExpired mengecek apakah token sudah kadaluarsa.
func (p *PasswordResetToken) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}

// IsUsed mengecek apakah token sudah pernah dipakai.
func (p *PasswordResetToken) IsUsed() bool {
	return p.UsedAt != nil
}
