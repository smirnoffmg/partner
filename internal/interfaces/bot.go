package interfaces

import (
	"time"

	"github.com/rs/zerolog/log"
	config "github.com/smirnoffmg/go-telegram-bot-template/config"
	usecases "github.com/smirnoffmg/go-telegram-bot-template/internal/usecases"
	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	Telebot *tele.Bot
}

func NewBot(cfg config.Config) (*Bot, error) {
	pref := tele.Settings{
		Token:  cfg.TelegramBotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot find bot with provided token")
		return nil, err
	}

	return &Bot{
		Telebot: b,
	}, nil
}

func (b *Bot) Start() {
	b.Telebot.Handle("/hello", func(c tele.Context) error {
		result := usecases.SayHello()
		return c.Send(result)
	})

	b.Telebot.Start()
}

func (b *Bot) Stop() {
	b.Telebot.Stop()
}
