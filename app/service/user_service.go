package service

import (
	"errors"
	"net/http"

	"keuangan-api/app/model"
	"keuangan-api/app/repository"

	"github.com/google/uuid"
)

type UserService struct {
	Repo *repository.UserRepository
}

func (s *UserService) GetRecentContacts(userID uuid.UUID) ([]model.User, int, error) {
	contacts, err := s.Repo.GetRecentContacts(userID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil kontak riwayat")
	}
	return contacts, http.StatusOK, nil
}
