package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type APIDocumentation struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Category        string          `gorm:"type:varchar(50);not null" json:"category"`
	Name            string          `gorm:"type:varchar(100);not null" json:"name"`
	Method          string          `gorm:"type:varchar(10);not null" json:"method"`
	Endpoint        string          `gorm:"type:varchar(255);not null" json:"endpoint"`
	IsProtected     bool            `gorm:"default:false" json:"is_protected"`
	Description     string          `gorm:"type:text" json:"description"`
	Headers         json.RawMessage `gorm:"type:jsonb" json:"headers"`
	RequestBody     json.RawMessage `gorm:"type:jsonb" json:"request_body"`
	ResponseSuccess json.RawMessage `gorm:"type:jsonb" json:"response_success"`
	ResponseError   json.RawMessage `gorm:"type:jsonb" json:"response_error"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func (APIDocumentation) TableName() string {
	return "api_documentations"
}
