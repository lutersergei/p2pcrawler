package repo

import (
	"github.com/lutersergei/p2pcrawler/pkg/chart"
	"github.com/lutersergei/p2pcrawler/pkg/price"
	"gorm.io/gorm"
)

type PriceRepo struct {
	db *gorm.DB
}

var query = `SELECT DISTINCT strftime('%s', created_at) / ? AS time,
    MIN(best_price) OVER w AS low,
    MAX(best_price) OVER w AS high,
    FIRST_VALUE(best_price) OVER w AS first,
    LAST_VALUE(best_price) OVER w AS last
FROM price_history
WHERE created_at >  datetime('now', '-' || ? || ' seconds')
WINDOW w AS (PARTITION BY strftime('%s', created_at) / ?);
`

func NewPriceRepo(db *gorm.DB) *PriceRepo {
	return &PriceRepo{db: db}
}

func (r *PriceRepo) Insert(history *price.PriceModel) error {
	result := r.db.Create(history)

	return result.Error
}

func (r *PriceRepo) HourChartData() []chart.HighLow {
	var data []chart.HighLow
	r.db.Raw(query, chart.MINUTE, chart.HOUR, chart.MINUTE).Scan(&data)

	return data
}

func (r *PriceRepo) DayChartData() []chart.HighLow {
	var data []chart.HighLow

	r.db.Raw(query, chart.MINUTE*30, chart.DAY, chart.MINUTE*30).Scan(&data)

	return data
}

func (r *PriceRepo) WeekChartData() []chart.HighLow {
	var data []chart.HighLow

	r.db.Raw(query, chart.HOUR*4, chart.WEEK, chart.HOUR*4).Scan(&data)

	return data
}

func (r *PriceRepo) MonthChartData() []chart.HighLow {
	var data []chart.HighLow

	r.db.Raw(query, chart.HOUR*24, chart.MONTH, chart.HOUR*24).Scan(&data)

	return data
}
