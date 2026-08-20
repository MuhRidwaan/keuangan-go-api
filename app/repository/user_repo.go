package repository

import (
	"keuangan-api/app/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.DB.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetRecentContacts mengambil daftar user unik yang pernah berada di grup tabungan/agenda yang sama dengan userID.
func (r *UserRepository) GetRecentContacts(userID uuid.UUID) ([]model.User, error) {
	var users []model.User
	query := `
		SELECT DISTINCT u.id, u.name, u.email, u.created_at, u.updated_at
		FROM users u
		WHERE u.id != ? AND u.deleted_at IS NULL AND (
			u.id IN (
				SELECT sm2.user_id FROM saving_members sm1
				JOIN saving_members sm2 ON sm1.goal_id = sm2.goal_id
				WHERE sm1.user_id = ?
			)
			OR
			u.id IN (
				SELECT am2.user_id FROM agenda_members am1
				JOIN agenda_members am2 ON am1.agenda_id = am2.agenda_id
				WHERE am1.user_id = ?
			)
		)
		ORDER BY u.name ASC
	`
	err := r.DB.Raw(query, userID, userID, userID).Scan(&users).Error
	return users, err
}
