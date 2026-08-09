package repository

import (
	"keuangan-api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgendaRepository struct {
	DB *gorm.DB
}

func (r *AgendaRepository) CreateAgendaWithOwner(agenda *model.Agenda, ownerID uuid.UUID, ownerEmail string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(agenda).Error; err != nil {
			return err
		}
		member := model.AgendaMember{
			AgendaID:  agenda.ID,
			UserID:    ownerID,
			Role:      model.RoleOwner,
			CreatedBy: &ownerEmail,
		}
		return tx.Create(&member).Error
	})
}

func (r *AgendaRepository) GetAgendasByMember(userID uuid.UUID) ([]model.Agenda, error) {
	var agendas []model.Agenda
	err := r.DB.
		Joins("JOIN agenda_members ON agenda_members.agenda_id = agendas.id").
		Where("agenda_members.user_id = ? AND agendas.deleted_at IS NULL", userID).
		Preload("Members").
		Order("agendas.start_date ASC").
		Find(&agendas).Error
	return agendas, err
}

func (r *AgendaRepository) FindByID(agendaID uuid.UUID) (*model.Agenda, error) {
	var agenda model.Agenda
	err := r.DB.Preload("Members").First(&agenda, "id = ?", agendaID).Error
	if err != nil {
		return nil, err
	}
	return &agenda, nil
}

func (r *AgendaRepository) GetMember(agendaID, userID uuid.UUID) (*model.AgendaMember, error) {
	var member model.AgendaMember
	err := r.DB.First(&member, "agenda_id = ? AND user_id = ?", agendaID, userID).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *AgendaRepository) AddMember(member *model.AgendaMember) error {
	return r.DB.Create(member).Error
}

func (r *AgendaRepository) UpdateAgenda(agenda *model.Agenda) error {
	return r.DB.Save(agenda).Error
}

// DeleteAgenda soft-delete agenda beserta semua member-nya dalam satu transaction.
func (r *AgendaRepository) DeleteAgenda(agendaID uuid.UUID) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Agenda{}, "id = ?", agendaID).Error; err != nil {
			return err
		}
		return tx.Delete(&model.AgendaMember{}, "agenda_id = ?", agendaID).Error
	})
}
