package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smirnoffmg/partner/internal/entities"
	"github.com/smirnoffmg/partner/internal/repositories"
)

type SubscriptionService struct {
	chatsRepo          repositories.IChatRepo
	invoicesRepo       repositories.IInvoiceRepo
	freeMessagesCount  int32
	pricePerMsgPack    int64
	msgPack            int32
	currency           string
	paymentDescription string
}

func NewSubscriptionService(chatsRepo repositories.IChatRepo, invoicesRepo repositories.IInvoiceRepo, freeMessagesCount int32) (*SubscriptionService, error) {
	return &SubscriptionService{
		chatsRepo:         chatsRepo,
		invoicesRepo:      invoicesRepo,
		freeMessagesCount: freeMessagesCount,
	}, nil
}

func (s *SubscriptionService) SetPaymentInfo(msgPack int32, pricePerMsgPack int64, currency string, paymentDescription string) {
	s.msgPack = msgPack
	s.pricePerMsgPack = pricePerMsgPack
	s.currency = currency
	s.paymentDescription = paymentDescription
}

func (s *SubscriptionService) IncreaseMessageCount(chatID int64) {
	if err := s.chatsRepo.IncreaseMessageCount(chatID); err != nil {
		log.Error().Err(err).Msgf("Cannot increase message count for chat with ID: %d", chatID)
	}
}

func (s *SubscriptionService) GetMessagesRemain(chatID int64) (int32, error) {
	chat, err := s.chatsRepo.GetOrCreate(chatID)
	if err != nil {
		log.Error().Err(err).Msgf("Cannot get chat with id: %d", chatID)
	}

	return s.freeMessagesCount + chat.PaidMessagesCount - chat.UserMessagesCount, nil
}

func (s *SubscriptionService) CreateInvoice(chatID int64) (*entities.Invoice, error) {
	invoice := &entities.Invoice{
		Title:          fmt.Sprintf("Invoice #%s", time.Now().Format("2006-01-02-15:04:05")),
		Description:    s.paymentDescription,
		Currency:       s.currency,
		Amount:         s.pricePerMsgPack,
		TelegramChatID: chatID,
	}

	if err := s.invoicesRepo.Create(invoice); err != nil {
		log.Error().Err(err).Msg("Cannot create invoice")

		return nil, err
	}

	return invoice, nil
}

func parsePayload(payload string) (int64, int64, error) {
	regexp := regexp.MustCompile(`^\d+#\d+$`)

	if !regexp.MatchString(payload) {
		return 0, 0, fmt.Errorf("invalid payload: %s", payload)
	}

	payloadParts := strings.Split(payload, "#")

	chatID, err := strconv.ParseInt(payloadParts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot parse chat ID: %s", payloadParts[0])
	}

	invoiceID, err := strconv.ParseInt(payloadParts[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot parse invoice ID: %s", payloadParts[1])
	}

	return chatID, invoiceID, nil
}

func (s *SubscriptionService) IsInvoiceOK(payload string) (bool, error) {
	_, invoiceID, err := parsePayload(payload)

	if err != nil {
		return false, err
	}

	invoice, err := s.invoicesRepo.Get(invoiceID)
	if err != nil {
		log.Error().Err(err).Msgf("Cannot get invoice with ID: %d", invoiceID)

		return false, err
	}

	return !invoice.Paid, nil
}

func (s *SubscriptionService) ProcessPayment(payload string) error {
	chatIDInt, invoiceIDInt, err := parsePayload(payload)

	if err != nil {
		return err
	}

	chat, err := s.chatsRepo.GetOrCreate(chatIDInt)

	if err != nil {
		return fmt.Errorf("cannot get chat with ID: %d", chatIDInt)
	}

	if _, err = s.invoicesRepo.Get(invoiceIDInt); err != nil {
		return fmt.Errorf("cannot get invoice with ID: %d", invoiceIDInt)
	}

	if err := s.invoicesRepo.Update(invoiceIDInt, map[string]interface{}{"paid": true}); err != nil {
		return fmt.Errorf("cannot update invoice with ID: %d", invoiceIDInt)
	}

	if err := s.chatsRepo.Update(chatIDInt, map[string]interface{}{"paid_messages_count": chat.PaidMessagesCount + s.msgPack}); err != nil {
		return fmt.Errorf("cannot update chat with ID: %d", chatIDInt)
	}

	return nil
}
