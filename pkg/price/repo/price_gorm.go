package repo

import (
	"gorm.io/gorm"
	"p2p_crawler/pkg/price"
)

type PriceRepo struct {
	db *gorm.DB
}

func NewPriceRepo(db *gorm.DB) *PriceRepo {
	return &PriceRepo{db: db}
}

func (r *PriceRepo) Insert(history *price.PriceHistory) error {
	result := r.db.Create(history)

	return result.Error
}
