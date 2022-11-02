package service

import (
	"bytes"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"image/color"
	"p2p_crawler/pkg/chart"
	"p2p_crawler/pkg/chart/gocandles"
)

type ChartService struct {
	logger *zap.SugaredLogger
	rep    PriceChartRepository
}

func NewChartService(logger *zap.SugaredLogger, rep PriceChartRepository) *ChartService {
	return &ChartService{logger: logger, rep: rep}
}

type PriceChartRepository interface {
	HourChartData() []chart.HighLow
	DayChartData() []chart.HighLow
	WeekChartData() []chart.HighLow
	MonthChartData() []chart.HighLow
}

func (svc *ChartService) Hour(ctx tele.Context) error {
	data := svc.rep.HourChartData()

	var candlesData []gocandles.Candle
	for i := range data {
		candlesData = append(candlesData, gocandles.Candle{
			Date:   int64(data[i].Time * chart.MINUTE),
			High:   data[i].High,
			Low:    data[i].Low,
			Open:   data[i].High,
			Close:  data[i].Low,
			Volume: 0,
		})
	}

	var buf bytes.Buffer

	err := gocandles.WriteChart(candlesData, getChartOption(), &buf)
	if err != nil {
		svc.logger.Error(err)
	}

	a := &tele.Photo{
		File: tele.FromReader(&buf),
	}

	return ctx.Send(a)
}

func (svc *ChartService) Day(ctx tele.Context) error {
	data := svc.rep.DayChartData()

	var candlesData []gocandles.Candle
	for i := range data {
		candlesData = append(candlesData, gocandles.Candle{
			Date:   int64(data[i].Time * chart.MINUTE * 30),
			High:   data[i].High,
			Low:    data[i].Low,
			Open:   data[i].High,
			Close:  data[i].Low,
			Volume: 0,
		})
	}

	var buf bytes.Buffer

	err := gocandles.WriteChart(candlesData, getChartOption(), &buf)
	if err != nil {
		svc.logger.Error(err)
	}

	a := &tele.Photo{
		File: tele.FromReader(&buf),
	}

	return ctx.Send(a)
}

func (svc *ChartService) Week(ctx tele.Context) error {
	data := svc.rep.WeekChartData()

	var candlesData []gocandles.Candle
	for i := range data {
		candlesData = append(candlesData, gocandles.Candle{
			Date:   int64(data[i].Time * chart.HOUR * 4),
			High:   data[i].High,
			Low:    data[i].Low,
			Open:   data[i].High,
			Close:  data[i].Low,
			Volume: 0,
		})
	}

	var buf bytes.Buffer

	err := gocandles.WriteChart(candlesData, getChartOption(), &buf)
	if err != nil {
		svc.logger.Error(err)
	}

	a := &tele.Photo{
		File: tele.FromReader(&buf),
	}

	return ctx.Send(a)
}

func (svc *ChartService) Month(ctx tele.Context) error {
	data := svc.rep.MonthChartData()

	var candlesData []gocandles.Candle
	for i := range data {
		candlesData = append(candlesData, gocandles.Candle{
			Date:   int64(data[i].Time * chart.DAY),
			High:   data[i].High,
			Low:    data[i].Low,
			Open:   data[i].High,
			Close:  data[i].Low,
			Volume: 0,
		})
	}

	var buf bytes.Buffer

	err := gocandles.WriteChart(candlesData, getChartOption(), &buf)
	if err != nil {
		svc.logger.Error(err)
	}

	a := &tele.Photo{
		File: tele.FromReader(&buf),
	}

	return ctx.Send(a)
}

func getChartOption() gocandles.Options {
	return gocandles.Options{
		LinesChartColor:      color.RGBA{0, 0, 0, 255},
		BackgroundChartColor: color.RGBA{255, 255, 255, 255},
		YLabelText:           "USDT",
		YLabelColor:          color.RGBA{255, 255, 255, 255},
		PositiveCandleColor:  color.RGBA{90, 90, 185, 255},
		NegativeCandleColor:  color.RGBA{255, 0, 0, 255},
		PikeCandleColor:      color.RGBA{211, 211, 211, 255},
		Width:                800,
		Height:               600,
		CandleWidth:          6,
		Rows:                 5,
		Columns:              7,
	}
}
