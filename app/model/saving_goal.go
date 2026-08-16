package model

import "time"

// SavingGoal merepresentasikan target tabungan bersama.
type SavingGoal struct {
	BaseModel
	Title         string    `gorm:"type:varchar(255);not null" json:"title"`
	TargetAmount  float64   `gorm:"type:numeric(15,2);not null" json:"target_amount"`
	CurrentAmount float64   `gorm:"type:numeric(15,2);not null;default:0" json:"current_amount"`
	Deadline      time.Time `gorm:"type:date;not null" json:"deadline"`

	// Relasi ke anggota dan kontribusi
	Members       []SavingMember       `gorm:"foreignKey:GoalID" json:"members,omitempty"`
	Contributions []SavingContribution `gorm:"foreignKey:GoalID" json:"contributions,omitempty"`
}
