package price

import (
	"time"
)

type PriceModel struct {
	ID            int       `json:"id" gorm:"primarykey"`
	BestPrice     float64   `json:"best_price"`
	Username      string    `json:"username"`
	Exchange      string    `json:"exchange"`
	RawJSON       string    `json:"rawJSON"`
	SurplusAmount float64   `json:"surplusAmount"`
	CreatedAt     time.Time `json:"created_at"`
}

type CurrentPriceResponse struct {
	ExchangeName  string
	BestPrice     float64
	SurplusAmount float64
}

func (PriceModel) TableName() string {
	return "price_history"
}

type ExchangeInterface interface {
	GetName() string
	DoRequest() (*PriceModel, error)
}

type PubSubInterface interface {
	//Subscribe(topic string) error
	//Unsubscribe(topic string) error
}
