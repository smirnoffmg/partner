package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	TelegramBotToken string `mapstructure:"TELEGRAM_BOT_TOKEN"` // this from env
	TelegramBotName  string `mapstructure:"telegram_bot_name"`  // this from config file
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot read config file")
		return
	}

	// this is a magic line that will allow viper to read environment variables
	viper.SetDefault("TELEGRAM_BOT_TOKEN", "")
	viper.AutomaticEnv()

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot unmarshal config")
		return
	}

	return
}
