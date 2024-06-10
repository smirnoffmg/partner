package repositories

import (
	"github.com/smirnoffmg/partner/internal/entities"

	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
)

type InvoicesRepo struct {
	db *gorm.DB
}

func NewInvoicesRepo(db *gorm.DB) *InvoicesRepo {
	return &InvoicesRepo{
		db: db,
	}
}

func (r *InvoicesRepo) Create(invoice *entities.Invoice) error {
	err := r.db.Create(invoice).Error
	if err != nil {
		log.Error().Err(err).Msgf("Cannot create invoice: %v", invoice)
	}

	return err
}

func (r *InvoicesRepo) Get(id int64) (*entities.Invoice, error) {
	var invoice entities.Invoice

	err := r.db.Model(&entities.Invoice{}).Where("id = ?", id).First(&invoice).Error

	if err != nil {
		log.Error().Err(err).Msgf("Cannot get invoice by id: %d", id)

		return nil, err
	}

	return &invoice, nil
}

func (r *InvoicesRepo) GetByChatID(id int64) (*[]entities.Invoice, error) {
	var invoices []entities.Invoice

	err := r.db.Where("id = ?", id).Find(&invoices).Error

	if err != nil {
		log.Error().Err(err).Msgf("Cannot get invoices by id: %d", id)

		return nil, err
	}

	return &invoices, nil
}

func (r *InvoicesRepo) Update(id int64, updates map[string]interface{}) error {
	err := r.db.Model(&entities.Invoice{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		log.Error().Err(err).Msgf("Cannot update invoice with id: %d", id)
	}

	return err
}
