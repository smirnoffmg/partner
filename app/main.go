package app

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	adapters "github.com/smirnoffmg/partner/internal/adapters"
	"github.com/smirnoffmg/partner/internal/ports"
	repo "github.com/smirnoffmg/partner/internal/repositories"
	"github.com/smirnoffmg/partner/internal/services"
)

type Config struct {
	TelegramBotToken  string `split_words:"true"`
	DBDsn             string `split_words:"true"`
	OpenaiAPIKey      string `split_words:"true"`
	OpenaiAssistantID string `split_words:"true"`
	Debug             bool   `split_words:"true"`
}

func loadConfig() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("PRTNR", &cfg); err != nil {
		return nil, err
	}

	log.Info().Interface("config", cfg).Msg("Configuration loaded")

	return &cfg, nil
}

type App struct {
	cfg *Config
	bot *ports.Bot
}

func NewApp() (*App, error) {
	log.Info().Msg("Loading configuration")

	cfg, err := loadConfig()
	if err != nil {
		log.Error().Err(err).Msg("Cannot load configuration")

		return nil, err
	}

	dbConn, err := adapters.NewDBConn(cfg.DBDsn)
	if err != nil {
		log.Error().Err(err).Msg("Cannot connect to database")

		return nil, err
	}

	chatsRepo := repo.NewChatsRepo(dbConn)

	chatGPTService, err := services.NewChatGPTService(chatsRepo, cfg.OpenaiAPIKey, cfg.OpenaiAssistantID)
	if err != nil {
		log.Error().Err(err).Msg("Cannot create chatGPT service")

		return nil, err
	}

	subscrService, err := services.NewSubscriptionService(chatsRepo)
	if err != nil {
		log.Error().Err(err).Msg("Cannot create subscription service")

		return nil, err
	}

	bot, err := ports.NewTGBot(cfg.TelegramBotToken, chatGPTService, subscrService)
	if err != nil {
		log.Error().Err(err).Msg("Cannot create bot")

		return nil, err
	}

	app := &App{
		cfg: cfg,
		bot: bot,
	}

	return app, nil
}

func (a *App) Start() error {
	return a.bot.Start()
}
