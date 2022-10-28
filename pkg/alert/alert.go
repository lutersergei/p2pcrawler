package alert

import "time"

const MoveUP MoveType = "up"
const MoveDown MoveType = "down"

const Sell DealType = "sell"
const Buy DealType = "buy"

const Active Status = "active"
const Done Status = "done"

type MoveType string
type DealType string
type Status string

type AlertDB struct {
	ID        int       `json:"id" gorm:"primarykey"`
	Price     float64   `json:"price"`
	Exchange  string    `json:"exchange"`
	Username  string    `json:"username"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	MoveType  MoveType  `json:"move_type"`
	DealType  DealType  `json:"deal_type"`
}

func (AlertDB) TableName() string {
	return "alert_request"
}
