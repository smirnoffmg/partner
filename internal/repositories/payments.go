package repositories

import (
	"github.com/rs/zerolog/log"
	"github.com/smirnoffmg/partner/internal/entities"

	"gorm.io/gorm"
)

type PaymentsRepo struct {
	db *gorm.DB
}

func NewPaymentsRepo(db *gorm.DB) *PaymentsRepo {
	return &PaymentsRepo{
		db: db,
	}
}

func (r *PaymentsRepo) Create(payment *entities.Payment) error {
	err := r.db.Create(payment).Error
	if err != nil {
		log.Error().Err(err).Msgf("Cannot create payment: %v", payment)
	}

	return err
}
