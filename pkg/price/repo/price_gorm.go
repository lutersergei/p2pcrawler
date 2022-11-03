package repo

import (
	"gorm.io/gorm"
	"p2p_crawler/pkg/chart"
	"p2p_crawler/pkg/price"
)

type PriceRepo struct {
	db *gorm.DB
}

func NewPriceRepo(db *gorm.DB) *PriceRepo {
	return &PriceRepo{db: db}
}

func (r *PriceRepo) Insert(history *price.PriceModel) error {
	result := r.db.Create(history)

	return result.Error
}

func (r *PriceRepo) HourChartData() []chart.HighLow {
	var data []chart.HighLow

	r.db.Raw("SELECT FLOOR(UNIX_TIMESTAMP(`created_at`)/(1 * 60)) AS time, " +
		"MAX(`best_price`) as high, MIN(`best_price`) as low FROM `price_history` " +
		"WHERE created_at >= date_sub(now(), INTERVAL 55 minute) " +
		"GROUP BY time").Scan(&data)

	return data
}

func (r *PriceRepo) DayChartData() []chart.HighLow {
	var data []chart.HighLow

	r.db.Raw("SELECT FLOOR(UNIX_TIMESTAMP(`created_at`)/(30 * 60)) AS time, " +
		"MAX(`best_price`) as high, MIN(`best_price`) as low FROM `price_history` " +
		"WHERE created_at >= date_sub(now(), INTERVAL 1 day) " +
		"GROUP BY time").Scan(&data)

	return data
}

func (r *PriceRepo) WeekChartData() []chart.HighLow {
	var data []chart.HighLow

	r.db.Raw("SELECT FLOOR(UNIX_TIMESTAMP(`created_at`)/(4 * 60 * 60)) AS time, " +
		"MAX(`best_price`) as high, MIN(`best_price`) as low FROM `price_history` " +
		"WHERE created_at >= date_sub(now(), INTERVAL 1 week) " +
		"GROUP BY time").Scan(&data)

	return data
}

func (r *PriceRepo) MonthChartData() []chart.HighLow {
	var data []chart.HighLow

	r.db.Raw("SELECT FLOOR(UNIX_TIMESTAMP(`created_at`)/(60 * 60 * 24)) AS time, " +
		"MAX(`best_price`) as high, MIN(`best_price`) as low FROM `price_history` " +
		"WHERE created_at >= date_sub(now(), INTERVAL 1 month) " +
		"GROUP BY time").Scan(&data)

	return data
}
