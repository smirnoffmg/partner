package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	config "github.com/smirnoffmg/go-telegram-bot-template/config"
	interfaces "github.com/smirnoffmg/go-telegram-bot-template/internal/interfaces"
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
	cfg, err := config.LoadConfig(".")

	if err != nil {
		log.Error().Err(err).Msg("Cannot load configuration")
		return err
	}

	bot, err := interfaces.NewBot(cfg)

	if err != nil {
		log.Error().Err(err).Msg("Cannot create bot")
		return err
	}

	go bot.Start()

	log.Info().Msg("Bot started")

	<-ctx.Done()

	log.Info().Msg("Shutting down (waiting 3 seconds)...")
	bot.Stop()

	<-time.After(3 * time.Second)

	return nil
}
