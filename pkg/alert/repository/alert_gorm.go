package repository

import (
	"gorm.io/gorm"
	"p2p_crawler/pkg/alert"
	"p2p_crawler/pkg/price"
)

type AlertRepo struct {
	db *gorm.DB
}

func NewAlertRepo(db *gorm.DB) *AlertRepo {
	return &AlertRepo{db: db}
}

func (a *AlertRepo) Match(history *price.PriceModel) []*alert.AlertDB {
	var alerts []*alert.AlertDB

	a.db.Where("price <= ? AND status = ?", history.BestPrice, alert.Active).Find(&alerts)

	return alerts
}

func (a *AlertRepo) Insert(model *alert.AlertDB) error {
	result := a.db.Omit("CreatedAt").Create(model)

	return result.Error
}

func (a *AlertRepo) Deactivate(model *alert.AlertDB) {
	model.Status = alert.Done
	a.db.Save(model)
}
