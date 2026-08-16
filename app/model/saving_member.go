package model

import (
	"time"

	"github.com/google/uuid"
)

// MemberRole mendefinisikan peran anggota dalam goal atau agenda.
type MemberRole string

const (
	RoleOwner  MemberRole = "owner"
	RoleMember MemberRole = "member"
)

// SavingMember adalah junction table antara SavingGoal dan User.
// Tidak menggunakan BaseModel penuh karena junction table tidak perlu soft delete & UUID PK sendiri.
type SavingMember struct {
	GoalID    uuid.UUID  `gorm:"type:uuid;primaryKey" json:"goal_id"`
	UserID    uuid.UUID  `gorm:"type:uuid;primaryKey" json:"user_id"`
	Role      MemberRole `gorm:"type:varchar(10);not null;default:'member'" json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy *string    `gorm:"type:varchar(255)" json:"created_by,omitempty"`

	// Relasi
	Goal SavingGoal `gorm:"foreignKey:GoalID" json:"goal,omitempty"`
	User User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
