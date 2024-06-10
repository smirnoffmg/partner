package adapters

import (
	"time"

	"github.com/rs/zerolog/log"
	entities "github.com/smirnoffmg/partner/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const retryTimeout = 5 * time.Second

func NewDBConn(dbDsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbDsn), &gorm.Config{})

	for err != nil {
		db, err = gorm.Open(postgres.Open(dbDsn), &gorm.Config{})

		<-time.After(retryTimeout)
	}

	log.Info().Msg("Connected to database")

	if err := db.AutoMigrate(&entities.Chat{}, &entities.Invoice{}); err != nil {
		log.Error().Err(err).Msg("Failed to migrate database")

		return nil, err
	}

	return db, nil
}
