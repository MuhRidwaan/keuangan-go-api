package service

import (
	"errors"
	"net/http"
	"time"

	"keuangan-api/app/model"
	"keuangan-api/app/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgendaService struct {
	Repo     *repository.AgendaRepository
	UserRepo *repository.UserRepository
}

// --- DTOs ---

type CreateAgendaInput struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     string  `json:"end_date" binding:"required"`
}

type UpdateAgendaInput struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
}

// AddAgendaMemberInput — invite by email
type AddAgendaMemberInput struct {
	Email string `json:"email" binding:"required,email"`
}

// --- Methods ---

func (s *AgendaService) CreateAgenda(input CreateAgendaInput, ownerID uuid.UUID, ownerEmail string) (*model.Agenda, int, error) {
	startDate, err := time.Parse(time.RFC3339, input.StartDate)
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("format start_date tidak valid, gunakan RFC3339 (contoh: 2006-01-02T15:04:05Z)")
	}
	endDate, err := time.Parse(time.RFC3339, input.EndDate)
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("format end_date tidak valid, gunakan RFC3339 (contoh: 2006-01-02T15:04:05Z)")
	}
	if !endDate.After(startDate) {
		return nil, http.StatusBadRequest, errors.New("end_date harus setelah start_date")
	}

	agenda := &model.Agenda{
		Title:       input.Title,
		Description: input.Description,
		StartDate:   startDate,
		EndDate:     endDate,
	}
	if err := s.Repo.CreateAgendaWithOwner(agenda, ownerID, ownerEmail); err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal membuat agenda")
	}
	return agenda, http.StatusCreated, nil
}

func (s *AgendaService) GetMyAgendas(userID uuid.UUID) ([]model.Agenda, int, error) {
	agendas, err := s.Repo.GetAgendasByMember(userID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil agenda")
	}
	return agendas, http.StatusOK, nil
}

func (s *AgendaService) Update(agendaID uuid.UUID, input UpdateAgendaInput, requesterID uuid.UUID) (*model.Agenda, int, error) {
	membership, err := s.Repo.GetMember(agendaID, requesterID)
	if err != nil || membership.Role != model.RoleOwner {
		return nil, http.StatusForbidden, errors.New("hanya owner yang dapat mengubah agenda ini")
	}

	agenda, err := s.Repo.FindByID(agendaID)
	if err != nil {
		return nil, http.StatusNotFound, errors.New("agenda tidak ditemukan")
	}

	if input.Title != "" {
		agenda.Title = input.Title
	}
	if input.Description != nil {
		agenda.Description = input.Description
	}

	// Update tanggal hanya jika keduanya diberikan, lalu validasi ulang
	newStart, newEnd := agenda.StartDate, agenda.EndDate
	if input.StartDate != "" {
		newStart, err = time.Parse(time.RFC3339, input.StartDate)
		if err != nil {
			return nil, http.StatusBadRequest, errors.New("format start_date tidak valid")
		}
	}
	if input.EndDate != "" {
		newEnd, err = time.Parse(time.RFC3339, input.EndDate)
		if err != nil {
			return nil, http.StatusBadRequest, errors.New("format end_date tidak valid")
		}
	}
	if !newEnd.After(newStart) {
		return nil, http.StatusBadRequest, errors.New("end_date harus setelah start_date")
	}
	agenda.StartDate = newStart
	agenda.EndDate = newEnd

	if err := s.Repo.UpdateAgenda(agenda); err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengupdate agenda")
	}
	return agenda, http.StatusOK, nil
}

func (s *AgendaService) Delete(agendaID uuid.UUID, requesterID uuid.UUID) (int, error) {
	membership, err := s.Repo.GetMember(agendaID, requesterID)
	if err != nil || membership.Role != model.RoleOwner {
		return http.StatusForbidden, errors.New("hanya owner yang dapat menghapus agenda ini")
	}

	if err := s.Repo.DeleteAgenda(agendaID); err != nil {
		return http.StatusInternalServerError, errors.New("gagal menghapus agenda")
	}
	return http.StatusOK, nil
}

func (s *AgendaService) AddMember(agendaID uuid.UUID, input AddAgendaMemberInput, requesterID uuid.UUID, requesterEmail string) (*model.AgendaMember, int, error) {
	requesterMembership, err := s.Repo.GetMember(agendaID, requesterID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusForbidden, errors.New("kamu bukan anggota agenda ini")
		}
		return nil, http.StatusInternalServerError, errors.New("gagal memvalidasi akses")
	}
	if requesterMembership.Role != model.RoleOwner {
		return nil, http.StatusForbidden, errors.New("hanya owner yang dapat menambah anggota")
	}

	// Lookup user by email
	targetUser, err := s.UserRepo.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.New("user dengan email tersebut tidak ditemukan")
		}
		return nil, http.StatusInternalServerError, errors.New("gagal mencari user")
	}

	// Cek duplikat
	_, err = s.Repo.GetMember(agendaID, targetUser.ID)
	if err == nil {
		return nil, http.StatusConflict, errors.New("user sudah menjadi anggota agenda ini")
	}

	member := &model.AgendaMember{
		AgendaID:  agendaID,
		UserID:    targetUser.ID,
		Role:      model.RoleMember,
		CreatedBy: &requesterEmail,
	}
	if err := s.Repo.AddMember(member); err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal menambah anggota")
	}
	return member, http.StatusCreated, nil
}
