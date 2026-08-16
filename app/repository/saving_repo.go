package repository

import (
	"keuangan-api/app/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SavingRepository struct {
	DB *gorm.DB
}

func (r *SavingRepository) CreateGoalWithOwner(goal *model.SavingGoal, ownerID uuid.UUID, ownerEmail string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(goal).Error; err != nil {
			return err
		}
		member := model.SavingMember{
			GoalID:    goal.ID,
			UserID:    ownerID,
			Role:      model.RoleOwner,
			CreatedBy: &ownerEmail,
		}
		return tx.Create(&member).Error
	})
}

func (r *SavingRepository) GetGoalsByMember(userID uuid.UUID) ([]model.SavingGoal, error) {
	var goals []model.SavingGoal
	err := r.DB.
		Joins("JOIN saving_members ON saving_members.goal_id = saving_goals.id").
		Where("saving_members.user_id = ? AND saving_goals.deleted_at IS NULL", userID).
		Preload("Members").
		Find(&goals).Error
	return goals, err
}

func (r *SavingRepository) FindGoalByID(goalID uuid.UUID) (*model.SavingGoal, error) {
	var goal model.SavingGoal
	err := r.DB.Preload("Members").First(&goal, "id = ?", goalID).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

func (r *SavingRepository) GetMember(goalID, userID uuid.UUID) (*model.SavingMember, error) {
	var member model.SavingMember
	err := r.DB.First(&member, "goal_id = ? AND user_id = ?", goalID, userID).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *SavingRepository) AddMember(member *model.SavingMember) error {
	return r.DB.Create(member).Error
}

// GetAllMembers mengambil semua member beserta data User-nya untuk keperluan notifikasi.
func (r *SavingRepository) GetAllMembers(goalID uuid.UUID) ([]model.SavingMember, error) {
	var members []model.SavingMember
	err := r.DB.Preload("User").Where("goal_id = ?", goalID).Find(&members).Error
	return members, err
}

// ContributionWithUser adalah DTO hasil join contributions + users.
type ContributionWithUser struct {
	model.SavingContribution
	ContributorName string `json:"contributor_name"`
}

// UserContributionSummary adalah ringkasan total kontribusi per user.
type UserContributionSummary struct {
	UserID   uuid.UUID `json:"user_id"`
	Name     string    `json:"name"`
	TotalIn  float64   `json:"total_in"`
	TotalOut float64   `json:"total_out"`
	Net      float64   `json:"net"`
}

// GetContributionHistory mengambil riwayat kontribusi beserta nama kontributor.
func (r *SavingRepository) GetContributionHistory(goalID uuid.UUID) ([]ContributionWithUser, error) {
	var results []ContributionWithUser
	err := r.DB.
		Table("saving_contributions sc").
		Select("sc.*, u.name AS contributor_name").
		Joins("JOIN users u ON u.id = sc.user_id").
		Where("sc.goal_id = ? AND sc.deleted_at IS NULL", goalID).
		Order("sc.date DESC").
		Scan(&results).Error
	return results, err
}

// GetContributionSummary mengambil ringkasan total per user untuk satu goal.
func (r *SavingRepository) GetContributionSummary(goalID uuid.UUID) ([]UserContributionSummary, error) {
	var summaries []UserContributionSummary
	err := r.DB.
		Table("saving_contributions sc").
		Select(`
			sc.user_id,
			u.name,
			COALESCE(SUM(CASE WHEN sc.type = 'in'  THEN sc.amount ELSE 0 END), 0) AS total_in,
			COALESCE(SUM(CASE WHEN sc.type = 'out' THEN sc.amount ELSE 0 END), 0) AS total_out,
			COALESCE(SUM(CASE WHEN sc.type = 'in'  THEN sc.amount ELSE -sc.amount END), 0) AS net
		`).
		Joins("JOIN users u ON u.id = sc.user_id").
		Where("sc.goal_id = ? AND sc.deleted_at IS NULL", goalID).
		Group("sc.user_id, u.name").
		Scan(&summaries).Error
	return summaries, err
}

func (r *SavingRepository) UpdateGoal(goal *model.SavingGoal) error {
	return r.DB.Save(goal).Error
}

// DeleteGoal soft-delete goal beserta semua member-nya dalam satu transaction.
func (r *SavingRepository) DeleteGoal(goalID uuid.UUID) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.SavingGoal{}, "id = ?", goalID).Error; err != nil {
			return err
		}
		return tx.Delete(&model.SavingMember{}, "goal_id = ?", goalID).Error
	})
}

// Contribute: insert contribution (type='in') + increment current_amount — atomic.
func (r *SavingRepository) Contribute(contribution *model.SavingContribution) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(contribution).Error; err != nil {
			return err
		}
		result := tx.Model(&model.SavingGoal{}).
			Where("id = ?", contribution.GoalID).
			UpdateColumn("current_amount", gorm.Expr("current_amount + ?", contribution.Amount))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

// Withdraw: insert contribution (type='out') + decrement current_amount — atomic.
// Caller wajib memvalidasi amount <= current_amount sebelum memanggil ini.
func (r *SavingRepository) Withdraw(contribution *model.SavingContribution) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(contribution).Error; err != nil {
			return err
		}
		result := tx.Model(&model.SavingGoal{}).
			Where("id = ?", contribution.GoalID).
			UpdateColumn("current_amount", gorm.Expr("current_amount - ?", contribution.Amount))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
