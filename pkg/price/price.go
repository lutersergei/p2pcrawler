package price

import (
	"time"
)

type PriceHistory struct {
	ID            int       `json:"id" gorm:"primarykey"`
	BestPrice     float64   `json:"best_price"`
	Username      string    `json:"username"`
	Exchange      string    `json:"exchange"`
	RawJSON       string    `json:"rawJSON"`
	SurplusAmount float64   `json:"surplusAmount"`
	CreatedAt     time.Time `json:"created_at"`
}

func (PriceHistory) TableName() string {
	return "price_history"
}
