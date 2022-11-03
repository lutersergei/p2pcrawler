package handler

import (
	"fmt"
	"github.com/lutersergei/p2pcrawler/pkg/alert"
	"github.com/lutersergei/p2pcrawler/pkg/config"
	"github.com/lutersergei/p2pcrawler/pkg/price"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type AlertTg struct {
	cfg    *config.Config
	bot    *tele.Bot
	logger *zap.SugaredLogger
}

func NewAlertTgHandler(cfg *config.Config, logger *zap.SugaredLogger, bot *tele.Bot) *AlertTg {

	return &AlertTg{cfg: cfg, bot: bot, logger: logger}
}

func (h *AlertTg) GetName() string {
	return "telegram"
}

func (h *AlertTg) Alert(al *alert.AlertDB, price *price.PriceModel) error {
	_, err := h.bot.Send(&tele.User{ID: int64(h.cfg.TgUser)}, fmt.Sprintf(
		"Price: %v, User: %s, Amount: %v",
		price.BestPrice,
		price.Username,
		price.SurplusAmount,
	))
	if err != nil {
		return fmt.Errorf("send to tg: %v", err)
	}

	return nil
}
