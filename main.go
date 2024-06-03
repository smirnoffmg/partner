package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smirnoffmg/partner/config"
	"github.com/smirnoffmg/partner/internal/adapters"
	"github.com/smirnoffmg/partner/internal/ports"
	repo "github.com/smirnoffmg/partner/internal/repositories"
	"github.com/smirnoffmg/partner/internal/services"
)

func main() {

	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Error().Err(err).Msg("Run")
		defer os.Exit(1)
	}

}

func run(ctx context.Context) error {
	log.Info().Msg("Loading configuration")
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Error().Err(err).Msg("Cannot load configuration")
		return err
	}

	dbConn, err := adapters.NewDBConn(cfg)

	if err != nil {
		log.Error().Err(err).Msg("Cannot connect to database")
		return err
	}

	chatsRepo := repo.NewChatsRepo(dbConn)

	chatGPTService, err := services.NewChatGPTService(chatsRepo, cfg.OpenaiApiKey, cfg.OpenaiAssistantId, cfg.Name)

	if err != nil {
		log.Error().Err(err).Msg("Cannot create chatGPTService")
		return err
	}

	bot, err := ports.NewTGBot(cfg.TelegramBotToken, chatGPTService)

	if err != nil {
		log.Error().Err(err).Msg("Cannot create bot")
		return err
	}

	if err := bot.Start(); err != nil {
		log.Error().Err(err).Msg("Problem inside bot.Start()")
		return err
	}

	<-ctx.Done()

	log.Info().Msg("Shutting down (waiting 3 seconds)...")

	<-time.After(3 * time.Second)

	return nil
}
