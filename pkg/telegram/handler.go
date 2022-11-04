package telegram

import (
	"fmt"
	"github.com/lutersergei/p2pcrawler/pkg/alert"
	al "github.com/lutersergei/p2pcrawler/pkg/alert/service"
	"github.com/lutersergei/p2pcrawler/pkg/chart/service"
	price "github.com/lutersergei/p2pcrawler/pkg/price/service"
	tele "gopkg.in/telebot.v3"
	"strconv"
	"time"
)

type TgHandler struct {
	navigationSvc *NavigationSvc
	priceSvc      *price.PriceService
	chartSvc      *service.ChartService
	alertSvc      *al.AlertService
}

func NewTgHandler(
	navigationSvc *NavigationSvc,
	priceSvc *price.PriceService,
	chartSvc *service.ChartService,
	alertSvc *al.AlertService,
) *TgHandler {
	return &TgHandler{
		navigationSvc: navigationSvc,
		priceSvc:      priceSvc,
		chartSvc:      chartSvc,
		alertSvc:      alertSvc,
	}
}

func (h *TgHandler) MainMenu(c tele.Context) error {
	return c.Send("Choose 👇", h.navigationSvc.MainMenu)
}

func (h *TgHandler) ChartMenu(c tele.Context) error {
	return c.Send("Choose 👇", h.navigationSvc.ChartMenu)
}

func (h *TgHandler) NotifyMenu(c tele.Context) error {
	return c.Send("Choose 👇", h.navigationSvc.NotifyMenu)
}

func (h *TgHandler) CurrentPrice(c tele.Context) error {
	var msg string

	r, err := h.priceSvc.CurrentPrice()
	if err != nil {
		return err
	}

	for _, response := range r {
		msg += fmt.Sprintf("%s: %v", response.ExchangeName, response.BestPrice)
	}

	return c.Send(msg)
}

func (h *TgHandler) AddNotification(c tele.Context) error {
	pr, err := strconv.ParseFloat(c.Message().Payload, 64)
	if err != nil {
		return fmt.Errorf("parsing payload: %v", err)
	}
	err = h.alertSvc.AddAlert(&alert.AlertDB{
		Price:     pr,
		Exchange:  "binance",
		Username:  c.Message().Sender.Username,
		Status:    alert.Active,
		CreatedAt: time.Time{},
		MoveType:  alert.MoveUP,
		DealType:  alert.Sell,
	})

	if err != nil {
		return fmt.Errorf("add alert: %v", err)
	}

	return c.Send(fmt.Sprintf("Add alert for: %v", c.Message().Payload))
}

func (h *TgHandler) Ping(c tele.Context) error {
	return c.Send("Pong!")
}

func (h *TgHandler) ChartHour(c tele.Context) error {
	return h.chartSvc.Hour(c)
}

func (h *TgHandler) ChartDay(c tele.Context) error {
	return h.chartSvc.Day(c)
}
func (h *TgHandler) ChartWeek(c tele.Context) error {
	return h.chartSvc.Week(c)
}
func (h *TgHandler) ChartMonth(c tele.Context) error {
	return h.chartSvc.Month(c)
}
