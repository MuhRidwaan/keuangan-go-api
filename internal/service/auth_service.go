package service

import (
	"errors"
	"net/http"

	"keuangan-api/internal/model"
	"keuangan-api/internal/repository"
	pkgjwt "keuangan-api/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	Repo *repository.AuthRepository
}

// RegisterInput adalah DTO untuk request register.
type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput adalah DTO untuk request login.
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterResult adalah data yang dikembalikan setelah register berhasil.
type RegisterResult struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// LoginResult adalah data yang dikembalikan setelah login berhasil.
type LoginResult struct {
	Token string `json:"token"`
}

func (s *AuthService) Register(input RegisterInput) (*RegisterResult, int, error) {
	// Cek apakah email sudah terdaftar
	_, err := s.Repo.FindByEmail(input.Email)
	if err == nil {
		return nil, http.StatusConflict, errors.New("email sudah terdaftar")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, http.StatusInternalServerError, errors.New("gagal memeriksa email")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal memproses password")
	}

	user := &model.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hashed),
	}

	if err := s.Repo.CreateUser(user); err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal menyimpan user")
	}

	return &RegisterResult{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}, http.StatusCreated, nil
}

func (s *AuthService) Login(input LoginInput) (*LoginResult, int, error) {
	user, err := s.Repo.FindByEmail(input.Email)
	if err != nil {
		// Jangan bocorkan apakah email ada atau tidak (security best practice)
		return nil, http.StatusUnauthorized, errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, http.StatusUnauthorized, errors.New("email atau password salah")
	}

	token, err := pkgjwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("gagal membuat token")
	}

	return &LoginResult{Token: token}, http.StatusOK, nil
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (s *AuthService) ChangePassword(userID uuid.UUID, input ChangePasswordInput) (int, error) {
	user, err := s.Repo.FindByID(userID.String())
	if err != nil {
		return http.StatusNotFound, errors.New("user tidak ditemukan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.OldPassword)); err != nil {
		return http.StatusBadRequest, errors.New("password lama tidak sesuai")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusInternalServerError, errors.New("gagal memproses password baru")
	}

	if err := s.Repo.UpdatePassword(user.ID.String(), string(hashed)); err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengupdate password")
	}

	return http.StatusOK, nil
}
