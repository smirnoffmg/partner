package repositories

import (
	"github.com/smirnoffmg/partner/internal/entities"
)

type IChatRepo interface {
	GetOrCreate(chatID int64) (*entities.Chat, error)
	Update(chatID int64, updates map[string]interface{}) error
	IncreaseMessageCount(chatID int64) error
}

type IInvoiceRepo interface {
	Create(invoice *entities.Invoice) error
	Get(id int64) (*entities.Invoice, error)
	GetByChatID(chatID int64) (*[]entities.Invoice, error)
	Update(id int64, updates map[string]interface{}) error
}

type IPaymentsRepo interface {
	Create(payment *entities.Payment) error
}
