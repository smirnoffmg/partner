package ports

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	svc "github.com/smirnoffmg/partner/internal/services"
)

type Bot struct {
	bot           *tgbotapi.BotAPI
	assistant     svc.IAssistant
	subscrService *svc.SubscriptionService
}

const (
	infoMessage       = "I'm a bot that can help you with your USMLE exam preparation.\nYou can contact my author @not_again_please for more info."
	startMessage      = "Hello! How shall we start?"
	freeMessagesCount = 50
)

func NewTGBot(tgBotToken string, assistant svc.IAssistant, subscrService *svc.SubscriptionService) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(tgBotToken)
	if err != nil {
		log.Error().Err(err).Msg("Cannot create bot")

		return nil, err
	}

	bot.Debug = false
	log.Info().Msgf("Authorized as %s", bot.Self.UserName)

	return &Bot{
		bot:           bot,
		assistant:     assistant,
		subscrService: subscrService,
	}, nil
}

func (b *Bot) handleCommand(update tgbotapi.Update) {
	if update.Message.Command() == "info" {
		infoMessageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, infoMessage)
		if _, err := b.bot.Send(infoMessageConfig); err != nil {
			log.Error().Err(err).Msg("Cannot send info message")
		}

		return
	}

	if update.Message.Command() == "start" {
		infoMessageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, startMessage)
		if _, err := b.bot.Send(infoMessageConfig); err != nil {
			log.Error().Err(err).Msg("Cannot send start message")
		}

		return
	}
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	var botMessage string

	log.Debug().Msgf("Message from %s", update.Message.From.UserName)

	answer, err := b.assistant.GetAnswer(update.Message.Chat.ID, update.Message.Text)
	if err != nil {
		log.Error().Err(err).Msg("Cannot get answer")

		botMessage = "Sorry, I can't answer right now. Please try again later."
	} else {
		b.subscrService.IncreaseMessageCount(update.Message.Chat.ID)

		botMessage = answer
	}

	answerMessageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, botMessage)
	answerMessageConfig.ReplyToMessageID = update.Message.MessageID

	if _, err := b.bot.Send(answerMessageConfig); err != nil {
		log.Error().Err(err).Msg("Cannot send answer")
	}
}

func (b *Bot) Start() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := b.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			go b.handleCommand(update)

			continue
		}

		go b.handleMessage(update)
	}

	return nil
}
