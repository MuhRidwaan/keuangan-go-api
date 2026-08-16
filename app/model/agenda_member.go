package model

import (
	"time"

	"github.com/google/uuid"
)

// AgendaMember adalah junction table antara Agenda dan User.
type AgendaMember struct {
	AgendaID  uuid.UUID  `gorm:"type:uuid;primaryKey" json:"agenda_id"`
	UserID    uuid.UUID  `gorm:"type:uuid;primaryKey" json:"user_id"`
	Role      MemberRole `gorm:"type:varchar(10);not null;default:'member'" json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy *string    `gorm:"type:varchar(255)" json:"created_by,omitempty"`

	// Relasi
	Agenda Agenda `gorm:"foreignKey:AgendaID" json:"agenda,omitempty"`
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
