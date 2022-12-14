package telegram

import (
	"fmt"
	"github.com/lutersergei/p2pcrawler/pkg/alert"
	al "github.com/lutersergei/p2pcrawler/pkg/alert/service"
	"github.com/lutersergei/p2pcrawler/pkg/chart/service"
	price "github.com/lutersergei/p2pcrawler/pkg/price/service"
	tele "gopkg.in/telebot.v3"
	"strconv"
	"strings"
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

func (h *TgHandler) GetAllNotify(c tele.Context) error {
	user := c.Sender().Username
	alerts, err := h.alertSvc.GetAlertsByUsername(user)
	if err != nil {
		return err
	}

	if len(alerts) == 0 {
		return c.Send("You haven't alerts")
	} else {
		var msg string
		for _, alertDB := range alerts {
			msg += fmt.Sprintf("Price: %v: Added: %v\n", alertDB.Price, alertDB.CreatedAt)
		}
		msg = strings.TrimSuffix(msg, "\n")

		return c.Send(msg)
	}
}

func (h *TgHandler) CurrentPrice(c tele.Context) error {
	var msg string

	r, err := h.priceSvc.CurrentPrice()
	if err != nil {
		return err
	}

	for _, response := range r {
		msg += fmt.Sprintf("%s: %v. Amount: %v", response.ExchangeName, response.BestPrice, response.SurplusAmount)
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
