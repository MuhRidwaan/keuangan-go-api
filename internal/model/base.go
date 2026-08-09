package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel adalah struct dasar yang di-embed oleh semua model.
// Menangani: UUID primary key, audit trail (created_by, updated_by, deleted_by), dan soft delete.
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy *string        `gorm:"type:varchar(255)" json:"created_by,omitempty"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy *string        `gorm:"type:varchar(255)" json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy *string        `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
}

// BeforeCreate hook: generate UUID jika belum ada
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
