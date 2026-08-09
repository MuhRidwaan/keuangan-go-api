package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"keuangan-api/internal/model"
	"keuangan-api/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PasswordResetService struct {
	Repo     *repository.PasswordResetRepository
	AuthRepo *repository.AuthRepository
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordInput struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// generateToken membuat random hex string 32 byte (64 karakter).
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *PasswordResetService) ForgotPassword(input ForgotPasswordInput) (int, error) {
	user, err := s.AuthRepo.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Jangan bocorkan apakah email terdaftar atau tidak (security best practice)
			// Tetap return 200 agar attacker tidak bisa enumerate email
			return http.StatusOK, nil
		}
		return http.StatusInternalServerError, errors.New("terjadi kesalahan pada server")
	}

	token, err := generateToken()
	if err != nil {
		return http.StatusInternalServerError, errors.New("gagal membuat token")
	}

	resetToken := &model.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Token berlaku 1 jam
	}

	if err := s.Repo.Create(resetToken); err != nil {
		return http.StatusInternalServerError, errors.New("gagal menyimpan token")
	}

	// MOCK: Cetak reset link ke console (ganti dengan email service di production)
	resetLink := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", token)
	fmt.Println("========================================")
	fmt.Printf("  [MOCK EMAIL] Kepada: %s\n", user.Email)
	fmt.Printf("  Subject: Reset Password Keuangan App\n")
	fmt.Printf("  Link: %s\n", resetLink)
	fmt.Printf("  Berlaku hingga: %s\n", resetToken.ExpiresAt.Format("02 Jan 2006 15:04:05"))
	fmt.Println("========================================")

	return http.StatusOK, nil
}

func (s *PasswordResetService) ResetPassword(input ResetPasswordInput) (int, error) {
	resetToken, err := s.Repo.FindByToken(input.Token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusBadRequest, errors.New("token tidak valid")
		}
		return http.StatusInternalServerError, errors.New("terjadi kesalahan pada server")
	}

	if resetToken.IsUsed() {
		return http.StatusBadRequest, errors.New("token sudah pernah digunakan")
	}

	if resetToken.IsExpired() {
		return http.StatusBadRequest, errors.New("token sudah kadaluarsa, silakan minta link baru")
	}

	// Hash password baru
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusInternalServerError, errors.New("gagal memproses password")
	}

	// Update password user
	if err := s.AuthRepo.UpdatePassword(resetToken.UserID.String(), string(hashed)); err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengupdate password")
	}

	// Tandai token sudah dipakai agar tidak bisa dipakai lagi
	if err := s.Repo.MarkAsUsed(resetToken.ID.String()); err != nil {
		return http.StatusInternalServerError, errors.New("gagal memvalidasi token")
	}

	return http.StatusOK, nil
}
