package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Name              string `split_words:"true"`
	TelegramBotToken  string `split_words:"true"`
	DbDsn             string `split_words:"true"`
	OpenaiApiKey      string `split_words:"true"`
	OpenaiAssistantId string `split_words:"true"`
	Debug             bool   `split_words:"true"`
}

func LoadConfig() (config Config, err error) {
	err = envconfig.Process("PRTNR", &config)

	if err != nil {
		log.Error().Err(err).Msg("Cannot load configuration")
		return
	}
	log.Info().Interface("config", config).Msg("Configuration loaded")
	return
}
