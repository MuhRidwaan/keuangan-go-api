package service

import (
	"errors"
	"net/http"

	"keuangan-api/internal/model"
	"keuangan-api/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationService struct {
	Repo *repository.NotificationRepository
}

func (s *NotificationService) GetMyNotifications(userID uuid.UUID) ([]model.Notification, int, error) {
	notifs, err := s.Repo.GetByUserID(userID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil notifikasi")
	}
	return notifs, http.StatusOK, nil
}

func (s *NotificationService) MarkAsRead(notifID uuid.UUID, userID uuid.UUID) (int, error) {
	err := s.Repo.MarkAsRead(notifID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, errors.New("notifikasi tidak ditemukan")
		}
		return http.StatusInternalServerError, errors.New("gagal mengupdate notifikasi")
	}
	return http.StatusOK, nil
}
