package ports

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	svc "github.com/smirnoffmg/partner/internal/services"
)

type Bot struct {
	bot           *tgbotapi.BotAPI
	assistant     svc.IAssistant
	subscrService *svc.SubscriptionService
	author        string
}

const (
	infoMessage       = "I'm a bot that can help you with your USMLE exam preparation.\nYou can contact my author %s for more info."
	startMessage      = "Hello! How shall we start?"
	payMessage        = "To continue please contact %s and tell him magic number %d. Price for 50 messages is 999 rubles."
	freeMessagesCount = 50
)

func NewTGBot(author string, tgBotToken string, assistant svc.IAssistant, subscrService *svc.SubscriptionService) (*Bot, error) {
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
		author:        author,
	}, nil
}

func (b *Bot) handleCommand(update tgbotapi.Update) {
	if update.Message.Command() == "info" {
		infoMsgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(infoMessage, b.author))
		if _, err := b.bot.Send(infoMsgConfig); err != nil {
			log.Error().Err(err).Msg("Cannot send info message")
		}

		return
	}

	if update.Message.Command() == "start" {
		startMsgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, startMessage)
		if _, err := b.bot.Send(startMsgConfig); err != nil {
			log.Error().Err(err).Msg("Cannot send start message")
		}

		return
	}

	if update.Message.Command() == "pay" {
		payMsgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(payMessage, b.author, update.Message.Chat.ID))
		if _, err := b.bot.Send(payMsgConfig); err != nil {
			log.Error().Err(err).Msg("Cannot send pay message")
		}

		return
	}
}

func (b *Bot) sendPayMessage(update tgbotapi.Update) {
	msg := "Sorry, you've ran out of messages. Please use /pay to continue"

	answerMessageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
	answerMessageConfig.ReplyToMessageID = update.Message.MessageID

	if _, err := b.bot.Send(answerMessageConfig); err != nil {
		log.Error().Err(err).Msg("Cannot send answer")
	}
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	var botMessage string

	chatID := update.Message.Chat.ID

	log.Debug().Msgf("Message from %s", update.Message.From.UserName)

	msgRemain, err := b.subscrService.GetMessagesRemain(chatID)
	if err != nil {
		log.Error().Err(err).Msg("Cannot get number of messages remain")
	} else if msgRemain < 1 {
		b.sendPayMessage(update)

		return
	}

	answer, err := b.assistant.GetAnswer(chatID, update.Message.Text)
	if err != nil {
		log.Error().Err(err).Msg("Cannot get answer")

		botMessage = "Sorry, I can't answer right now. Please try again later."
	} else {
		b.subscrService.IncreaseMessageCount(chatID)

		botMessage = answer
	}

	botMessage = fmt.Sprintf("%s\n\n(messages left: %d)", botMessage, msgRemain)

	answerMessageConfig := tgbotapi.NewMessage(chatID, botMessage)
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
