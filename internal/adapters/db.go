package adapters

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
	config "github.com/smirnoffmg/partner/config"
	entities "github.com/smirnoffmg/partner/internal/entities"
)

func NewDBConn(cfg config.Config) (*gorm.DB, error) {
	// wait for the database to be ready

	db, err := gorm.Open(postgres.Open(cfg.DbDsn), &gorm.Config{})

	for err != nil {
		db, err = gorm.Open(postgres.Open(cfg.DbDsn), &gorm.Config{})
		<-time.After(5 * time.Second)
	}

	log.Info().Msg("Connected to database")

	if err := db.AutoMigrate(&entities.Chat{}); err != nil {
		log.Error().Err(err).Msg("Failed to migrate database")
		return nil, err
	}

	return db, nil
}
