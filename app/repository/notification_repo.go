package repository

import (
	"keuangan-api/app/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	DB *gorm.DB
}

// Create menyimpan satu notifikasi baru.
func (r *NotificationRepository) Create(notification *model.Notification) error {
	return r.DB.Create(notification).Error
}

// BulkCreate menyimpan banyak notifikasi sekaligus dalam satu query.
func (r *NotificationRepository) BulkCreate(notifications []model.Notification) error {
	if len(notifications) == 0 {
		return nil
	}
	return r.DB.Create(&notifications).Error
}

// GetByUserID mengambil semua notifikasi milik user, terbaru dulu.
func (r *NotificationRepository) GetByUserID(userID uuid.UUID) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

// MarkAsRead menandai satu notifikasi sebagai sudah dibaca.
// Hanya berhasil jika notifikasi milik userID yang diberikan (ownership check).
func (r *NotificationRepository) MarkAsRead(notifID, userID uuid.UUID) error {
	result := r.DB.Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", notifID, userID).
		Update("is_read", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
