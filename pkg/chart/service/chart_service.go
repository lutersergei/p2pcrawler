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

	candlesData := buildCandles(data, chart.MINUTE)

	chrt, err := buildChart(candlesData)
	if err != nil {
		return err
	}

	return ctx.Send(chrt)
}

func (svc *ChartService) Day(ctx tele.Context) error {
	data := svc.rep.DayChartData()

	candlesData := buildCandles(data, chart.MINUTE*30)

	chrt, err := buildChart(candlesData)
	if err != nil {
		return err
	}

	return ctx.Send(chrt)
}

func (svc *ChartService) Week(ctx tele.Context) error {
	data := svc.rep.WeekChartData()

	candlesData := buildCandles(data, chart.HOUR*4)

	chrt, err := buildChart(candlesData)
	if err != nil {
		return err
	}

	return ctx.Send(chrt)
}

func (svc *ChartService) Month(ctx tele.Context) error {
	data := svc.rep.MonthChartData()

	candlesData := buildCandles(data, chart.DAY)

	chrt, err := buildChart(candlesData)
	if err != nil {
		return err
	}

	return ctx.Send(chrt)
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

func buildCandles(data []chart.HighLow, candleDur int) []gocandles.Candle {
	var prevHigh, prevLow float64
	var candlesData []gocandles.Candle
	for i := range data {
		if (prevHigh > data[i].High && prevLow > data[i].Low) || (prevHigh == data[i].High && prevLow > data[i].Low) {
			candlesData = append(candlesData, gocandles.Candle{
				Date:  int64(data[i].Time * candleDur),
				High:  data[i].High,
				Low:   data[i].Low,
				Open:  data[i].High,
				Close: data[i].Low,
			})
		} else {
			candlesData = append(candlesData, gocandles.Candle{
				Date:  int64(data[i].Time * candleDur),
				High:  data[i].High,
				Low:   data[i].Low,
				Open:  data[i].Low,
				Close: data[i].High,
			})
		}
		prevHigh = data[i].High
		prevLow = data[i].Low
	}

	return candlesData
}

func buildChart(data []gocandles.Candle) (tele.Sendable, error) {
	var buf bytes.Buffer

	err := gocandles.WriteChart(data, getChartOption(), &buf)
	if err != nil {
		return nil, err
	}

	a := &tele.Photo{
		File: tele.FromReader(&buf),
	}

	return a, nil
}
