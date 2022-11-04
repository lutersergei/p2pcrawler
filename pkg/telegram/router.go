package telegram

import tele "gopkg.in/telebot.v3"

type TgRouter struct {
	handler *TgHandler
	bot     *tele.Bot
	nav     *NavigationSvc
}

func NewTgRouter(handler *TgHandler, bot *tele.Bot, nav *NavigationSvc) *TgRouter {
	return &TgRouter{handler: handler, bot: bot, nav: nav}
}

func (r *TgRouter) ApplyRoutes() {
	r.bot.Handle("/start", r.handler.MainMenu)
	r.bot.Handle(&r.nav.BtnHome, r.handler.MainMenu)
	r.bot.Handle(&r.nav.BtnChart, r.handler.ChartMenu)
	r.bot.Handle(&r.nav.BtnNotify, r.handler.NotifyMenu)

	r.bot.Handle(&r.nav.BtnCurrentPrice, r.handler.CurrentPrice)
	r.bot.Handle(&r.nav.BtnPing, r.handler.Ping)

	r.bot.Handle(&r.nav.BtnChartHour, r.handler.ChartHour)
	r.bot.Handle(&r.nav.BtnChartDay, r.handler.ChartDay)
	r.bot.Handle(&r.nav.BtnChartWeek, r.handler.ChartWeek)
	r.bot.Handle(&r.nav.BtnChartMonth, r.handler.ChartMonth)

	r.bot.Handle("/add", r.handler.AddNotification)
}
