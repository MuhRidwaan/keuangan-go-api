package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"keuangan-api/internal/model"
	"keuangan-api/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SavingService struct {
	Repo     *repository.SavingRepository
	UserRepo *repository.UserRepository
	NotifRepo *repository.NotificationRepository
}

// --- DTOs ---

type CreateGoalInput struct {
	Title        string  `json:"title" binding:"required"`
	TargetAmount float64 `json:"target_amount" binding:"required,gt=0"`
	Deadline     string  `json:"deadline" binding:"required"`
}

type UpdateGoalInput struct {
	Title        string  `json:"title"`
	TargetAmount float64 `json:"target_amount" binding:"omitempty,gt=0"`
	Deadline     string  `json:"deadline"`
}

// AddMemberByEmailInput — invite by email, bukan user_id
type AddMemberByEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}

type ContributeInput struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Date   string  `json:"date" binding:"required"`
}

type WithdrawInput struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Date   string  `json:"date" binding:"required"`
}

// --- Methods ---

func (s *SavingService) CreateGoal(input CreateGoalInput, ownerID uuid.UUID, ownerEmail string) (*model.SavingGoal, int, error) {
	deadline, err := time.Parse("2006-01-02", input.Deadline)
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("format deadline tidak valid, gunakan YYYY-MM-DD")
	}

	goal := &model.SavingGoal{
		Title:        input.Title,
		TargetAmount: input.TargetAmount,
		Deadline:     deadline,
	}
	if err := s.Repo.CreateGoalWithOwner(goal, ownerID, ownerEmail); err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal membuat saving goal")
	}
	return goal, http.StatusCreated, nil
}

func (s *SavingService) GetMyGoals(userID uuid.UUID) ([]model.SavingGoal, int, error) {
	goals, err := s.Repo.GetGoalsByMember(userID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil saving goals")
	}
	return goals, http.StatusOK, nil
}

func (s *SavingService) Update(goalID uuid.UUID, input UpdateGoalInput, requesterID uuid.UUID) (*model.SavingGoal, int, error) {
	// Hanya owner yang boleh edit
	membership, err := s.Repo.GetMember(goalID, requesterID)
	if err != nil || membership.Role != model.RoleOwner {
		return nil, http.StatusForbidden, errors.New("hanya owner yang dapat mengubah goal ini")
	}

	goal, err := s.Repo.FindGoalByID(goalID)
	if err != nil {
		return nil, http.StatusNotFound, errors.New("saving goal tidak ditemukan")
	}

	if input.Title != "" {
		goal.Title = input.Title
	}
	if input.TargetAmount > 0 {
		goal.TargetAmount = input.TargetAmount
	}
	if input.Deadline != "" {
		deadline, err := time.Parse("2006-01-02", input.Deadline)
		if err != nil {
			return nil, http.StatusBadRequest, errors.New("format deadline tidak valid, gunakan YYYY-MM-DD")
		}
		goal.Deadline = deadline
	}

	if err := s.Repo.UpdateGoal(goal); err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengupdate goal")
	}
	return goal, http.StatusOK, nil
}

func (s *SavingService) Delete(goalID uuid.UUID, requesterID uuid.UUID) (int, error) {
	membership, err := s.Repo.GetMember(goalID, requesterID)
	if err != nil || membership.Role != model.RoleOwner {
		return http.StatusForbidden, errors.New("hanya owner yang dapat menghapus goal ini")
	}

	if err := s.Repo.DeleteGoal(goalID); err != nil {
		return http.StatusInternalServerError, errors.New("gagal menghapus goal")
	}
	return http.StatusOK, nil
}

func (s *SavingService) AddMember(goalID uuid.UUID, input AddMemberByEmailInput, requesterID uuid.UUID, requesterEmail string) (*model.SavingMember, int, error) {
	// Validasi: hanya owner
	requesterMembership, err := s.Repo.GetMember(goalID, requesterID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusForbidden, errors.New("kamu bukan anggota goal ini")
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
	_, err = s.Repo.GetMember(goalID, targetUser.ID)
	if err == nil {
		return nil, http.StatusConflict, errors.New("user sudah menjadi anggota goal ini")
	}

	member := &model.SavingMember{
		GoalID:    goalID,
		UserID:    targetUser.ID,
		Role:      model.RoleMember,
		CreatedBy: &requesterEmail,
	}
	if err := s.Repo.AddMember(member); err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal menambah anggota")
	}
	return member, http.StatusCreated, nil
}

// broadcastNotification mengirim notifikasi ke semua member goal KECUALI actor itu sendiri.
// Dipanggil secara fire-and-forget — error tidak menggagalkan operasi utama.
func (s *SavingService) broadcastNotification(goalID uuid.UUID, actorID uuid.UUID, title, message string) {
	members, err := s.Repo.GetAllMembers(goalID)
	if err != nil {
		return
	}

	var notifications []model.Notification
	for _, m := range members {
		if m.UserID == actorID {
			continue // skip actor, notif hanya untuk member lain
		}
		notifications = append(notifications, model.Notification{
			UserID:  m.UserID,
			Title:   title,
			Message: message,
		})
	}

	// Best-effort: error diabaikan agar tidak menggagalkan transaksi utama
	_ = s.NotifRepo.BulkCreate(notifications)
}

func (s *SavingService) Contribute(goalID uuid.UUID, input ContributeInput, userID uuid.UUID) (*model.SavingContribution, int, error) {
	_, err := s.Repo.GetMember(goalID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusForbidden, errors.New("kamu bukan anggota goal ini")
		}
		return nil, http.StatusInternalServerError, errors.New("gagal memvalidasi akses")
	}

	// Ambil goal untuk nama di notifikasi
	goal, err := s.Repo.FindGoalByID(goalID)
	if err != nil {
		return nil, http.StatusNotFound, errors.New("saving goal tidak ditemukan")
	}

	// Ambil data actor untuk nama di notifikasi
	actor, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil data user")
	}

	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("format tanggal tidak valid, gunakan YYYY-MM-DD")
	}

	contribution := &model.SavingContribution{
		GoalID: goalID,
		UserID: userID,
		Amount: input.Amount,
		Type:   model.ContributionIn,
		Date:   date,
	}
	if err := s.Repo.Contribute(contribution); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.New("saving goal tidak ditemukan")
		}
		return nil, http.StatusInternalServerError, errors.New("gagal menyimpan kontribusi")
	}

	// Broadcast notifikasi ke semua member lain (fire-and-forget)
	go s.broadcastNotification(
		goalID, userID,
		"Kontribusi Baru",
		fmt.Sprintf("%s baru saja menyetor Rp %.0f ke tabungan \"%s\"", actor.Name, input.Amount, goal.Title),
	)

	return contribution, http.StatusCreated, nil
}

func (s *SavingService) Withdraw(goalID uuid.UUID, input WithdrawInput, userID uuid.UUID) (*model.SavingContribution, int, error) {
	_, err := s.Repo.GetMember(goalID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusForbidden, errors.New("kamu bukan anggota goal ini")
		}
		return nil, http.StatusInternalServerError, errors.New("gagal memvalidasi akses")
	}

	goal, err := s.Repo.FindGoalByID(goalID)
	if err != nil {
		return nil, http.StatusNotFound, errors.New("saving goal tidak ditemukan")
	}

	if input.Amount > goal.CurrentAmount {
		return nil, http.StatusBadRequest, errors.New("jumlah penarikan melebihi saldo saat ini")
	}

	actor, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil data user")
	}

	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("format tanggal tidak valid, gunakan YYYY-MM-DD")
	}

	contribution := &model.SavingContribution{
		GoalID: goalID,
		UserID: userID,
		Amount: input.Amount,
		Type:   model.ContributionOut,
		Date:   date,
	}
	if err := s.Repo.Withdraw(contribution); err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal memproses penarikan")
	}

	// Broadcast notifikasi ke semua member lain (fire-and-forget)
	go s.broadcastNotification(
		goalID, userID,
		"Penarikan Dana",
		fmt.Sprintf("%s baru saja menarik Rp %.0f dari tabungan \"%s\"", actor.Name, input.Amount, goal.Title),
	)

	return contribution, http.StatusCreated, nil
}

// ContributionHistoryResult adalah response untuk GET /savings/:id/contributions.
type ContributionHistoryResult struct {
	History []repository.ContributionWithUser    `json:"history"`
	Summary []repository.UserContributionSummary `json:"summary"`
}

func (s *SavingService) GetContributionHistory(goalID uuid.UUID, requesterID uuid.UUID) (*ContributionHistoryResult, int, error) {
	// Validasi: hanya member yang boleh lihat history
	_, err := s.Repo.GetMember(goalID, requesterID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusForbidden, errors.New("kamu bukan anggota goal ini")
		}
		return nil, http.StatusInternalServerError, errors.New("gagal memvalidasi akses")
	}

	history, err := s.Repo.GetContributionHistory(goalID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil riwayat kontribusi")
	}

	summary, err := s.Repo.GetContributionSummary(goalID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal mengambil ringkasan kontribusi")
	}

	// Pastikan slice tidak nil agar response JSON selalu array, bukan null
	if history == nil {
		history = []repository.ContributionWithUser{}
	}
	if summary == nil {
		summary = []repository.UserContributionSummary{}
	}

	return &ContributionHistoryResult{
		History: history,
		Summary: summary,
	}, http.StatusOK, nil
}
