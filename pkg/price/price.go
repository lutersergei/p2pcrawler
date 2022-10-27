package price

import (
	"time"
)

type PriceHistoryDB struct {
	ID            int       `json:"id" gorm:"primarykey"`
	MaxPrice      float64   `json:"maxPrice"`
	Username      string    `json:"username"`
	Exchange      string    `json:"exchange"`
	RawJSON       string    `json:"rawJSON"`
	SurplusAmount float64   `json:"surplusAmount"`
	CreatedAt     time.Time `json:"created_at"`
}

func (PriceHistoryDB) TableName() string {
	return "price_history"
}
