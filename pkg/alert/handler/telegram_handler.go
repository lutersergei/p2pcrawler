package handler

import (
	"fmt"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"p2p_crawler/pkg/alert"
	"p2p_crawler/pkg/config"
	"p2p_crawler/pkg/price"
	"time"
)

type TelegramHandler struct {
	cfg    *config.Config
	bot    *tele.Bot
	logger *zap.SugaredLogger
}

func NewTelegramHandler(cfg *config.Config, logger *zap.SugaredLogger) *TelegramHandler {
	// bot settings
	pref := tele.Settings{
		Token:  cfg.TgBotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		logger.Fatalw("connect to bot: ", zap.Error(err))
		return nil
	}
	return &TelegramHandler{cfg: cfg, bot: bot, logger: logger}
}

func (t *TelegramHandler) GetName() string {
	return "telegram"
}

func (t *TelegramHandler) Alert(al *alert.AlertDB, price *price.PriceHistory) error {
	_, err := t.bot.Send(&tele.User{ID: int64(t.cfg.TgUser)}, fmt.Sprintf(
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
