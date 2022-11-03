package handler

import (
	"fmt"
	price "github.com/lutersergei/p2pcrawler/pkg/price/service"
	tele "gopkg.in/telebot.v3"
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
}
