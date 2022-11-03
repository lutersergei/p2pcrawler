package handler

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	price "p2p_crawler/pkg/price/service"
)

type PriceTg struct {
	priceSvc *price.PriceService
}

func NewPriceTgHandler(priceSvc *price.PriceService) *PriceTg {
	return &PriceTg{priceSvc: priceSvc}
}

func (h *PriceTg) CurrentPrice(c tele.Context) error {
	var msg string

	r, err := h.priceSvc.CurrentPrice()
	if err != nil {
		return err
	}

	for _, response := range r {
		msg += fmt.Sprintf("%s: %v", response.ExchangeName, response.BestPrice)
	}

	return c.Send(msg)
	//_, err = h.bot.Send(&tele.User{ID: int64(h.cfg.TgUser)}, msg)
	//if err != nil {
	//	return fmt.Errorf("send to tg: %v", err)
	//}
	//
	//return nil
}
