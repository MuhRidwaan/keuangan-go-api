package model

// User merepresentasikan akun pengguna aplikasi.
type User struct {
	BaseModel
	Name         string `gorm:"type:varchar(255);not null" json:"name"`
	Email        string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"` // json:"-" agar tidak pernah terekspos ke response
}
