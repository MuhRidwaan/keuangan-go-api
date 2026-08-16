package model

import "time"

// Agenda merepresentasikan acara atau jadwal bersama.
type Agenda struct {
	BaseModel
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	StartDate   time.Time `gorm:"type:timestamptz;not null" json:"start_date"`
	EndDate     time.Time `gorm:"type:timestamptz;not null" json:"end_date"`

	// Relasi ke anggota agenda
	Members []AgendaMember `gorm:"foreignKey:AgendaID" json:"members,omitempty"`
}
