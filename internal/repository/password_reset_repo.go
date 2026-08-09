package repository

import (
	"time"

	"keuangan-api/internal/model"

	"gorm.io/gorm"
)

type PasswordResetRepository struct {
	DB *gorm.DB
}

// Create menyimpan token baru, setelah menghapus token lama milik user yang sama.
func (r *PasswordResetRepository) Create(token *model.PasswordResetToken) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Hapus token lama milik user ini agar tidak menumpuk
		if err := tx.Where("user_id = ?", token.UserID).Delete(&model.PasswordResetToken{}).Error; err != nil {
			return err
		}
		return tx.Create(token).Error
	})
}

// FindByToken mencari token yang valid (belum expired, belum dipakai).
func (r *PasswordResetRepository) FindByToken(tokenStr string) (*model.PasswordResetToken, error) {
	var token model.PasswordResetToken
	err := r.DB.Where("token = ?", tokenStr).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// MarkAsUsed menandai token sudah dipakai dengan set used_at = now().
func (r *PasswordResetRepository) MarkAsUsed(tokenID string) error {
	now := time.Now()
	return r.DB.Model(&model.PasswordResetToken{}).
		Where("id = ?", tokenID).
		Update("used_at", now).Error
}
