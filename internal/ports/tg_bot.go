package ports

import (
	"encoding/json"
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
	paymentToken  string
}

const (
	infoMessage       = "I'm a bot that can help you with your USMLE exam preparation.\nYou can contact my author %s for more info."
	startMessage      = "Hello! How shall we start?"
	payMessage        = "To continue please contact %s and tell him magic number %d. Price for 50 messages is 999 rubles."
	freeMessagesCount = 50
)

func NewTGBot(author string, tgBotToken string, assistant svc.IAssistant, subscrService *svc.SubscriptionService, paymentToken string) (*Bot, error) {
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
		paymentToken:  paymentToken,
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
		// create and send invoice
		invoice, err := b.subscrService.CreateInvoice(update.Message.Chat.ID)

		if err != nil {
			log.Error().Err(err).Msg("Cannot create invoice")

			return
		}

		payload := fmt.Sprintf("%d#%d", invoice.TelegramChatID, invoice.ID)

		prices := []tgbotapi.LabeledPrice{
			{
				Label:  invoice.Description,
				Amount: int(invoice.Amount),
			},
		}

		invoiceMsg := tgbotapi.NewInvoice(
			invoice.TelegramChatID,
			invoice.Title,
			invoice.Description,
			payload,
			b.paymentToken,
			"",
			invoice.Currency,
			prices)

		invoiceMsg.MaxTipAmount = 0
		invoiceMsg.SuggestedTipAmounts = []int{}
		invoiceMsg.NeedEmail = true
		invoiceMsg.SendEmailToProvider = true

		providerData := map[string]interface{}{
			"receipt": map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"description": invoice.Description,
						"amount": map[string]string{
							"value":    fmt.Sprintf("%d.00", invoice.Amount/100),
							"currency": invoice.Currency,
						},
						"vat_code": 1,
						"quantity": 1,
					},
				},
			},
		}

		var providerDataJSON []byte

		providerDataJSON, err = json.Marshal(providerData)

		if err != nil {
			log.Error().Err(err).Msg("Cannot marshal providerData")

			return
		}

		log.Debug().Interface("providerData", string(providerDataJSON)).Msg("ProviderData")

		invoiceMsg.ProviderData = string(providerDataJSON)

		if _, err := b.bot.Send(invoiceMsg); err != nil {
			log.Error().Err(err).Msg("Cannot send invoice")
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

func (b *Bot) handlePreCheckoutQuery(update tgbotapi.Update) {
	log.Info().Interface("preCheckoutQuery", update.PreCheckoutQuery).Msg("PreCheckoutQuery")

	errMsg := ""
	check, err := b.subscrService.IsInvoiceOK(update.PreCheckoutQuery.InvoicePayload)

	if err != nil {
		log.Error().Err(err).Msg("Cannot check invoice")

		errMsg = "Cannot check invoice"
	}

	if !check {
		errMsg = "Invoice already paid"
	}

	pca := tgbotapi.PreCheckoutConfig{
		PreCheckoutQueryID: update.PreCheckoutQuery.ID,
		OK:                 check,
		ErrorMessage:       errMsg,
	}

	if _, err := b.bot.Request(pca); err != nil {
		log.Error().Err(err).Msg("Cannot answer preCheckoutQuery")
	}
}

func (b *Bot) handlePayment(update tgbotapi.Update) {
	log.Info().Interface("payment", update.Message.SuccessfulPayment).Msg("SuccessfulPayment")

	err := b.subscrService.ProcessPayment(update.Message.SuccessfulPayment.InvoicePayload)

	if err != nil {
		log.Error().Err(err).Msg("Cannot process payment")

		return
	}

	newMessagesRemain, err := b.subscrService.GetMessagesRemain(update.Message.Chat.ID)

	if err != nil {
		log.Error().Err(err).Msg("Cannot get messages remain")

		return
	}

	msg := fmt.Sprintf("Thank you for payment! You have %d messages left", newMessagesRemain)

	answerMessageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, msg)

	if _, err := b.bot.Send(answerMessageConfig); err != nil {
		log.Error().Err(err).Msg("Cannot send answer")
	}
}

func (b *Bot) Start() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := b.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.PreCheckoutQuery != nil {
			go b.handlePreCheckoutQuery(update)

			continue
		}

		if update.Message.SuccessfulPayment != nil {
			go b.handlePayment(update)

			continue
		}

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
