package model

import (
	"time"

	"github.com/google/uuid"
)

// ContributionType membedakan setoran masuk dan penarikan.
type ContributionType string

const (
	ContributionIn  ContributionType = "in"
	ContributionOut ContributionType = "out"
)

// SavingContribution merepresentasikan satu setoran atau penarikan dari SavingGoal.
type SavingContribution struct {
	BaseModel
	GoalID uuid.UUID        `gorm:"type:uuid;not null;index" json:"goal_id"`
	Goal   SavingGoal       `gorm:"foreignKey:GoalID" json:"goal,omitempty"`
	UserID uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	User   User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Amount float64          `gorm:"type:numeric(15,2);not null" json:"amount"`
	Type   ContributionType `gorm:"type:varchar(3);not null;default:'in'" json:"type"`
	Date   time.Time        `gorm:"type:date;not null" json:"date"`
}
