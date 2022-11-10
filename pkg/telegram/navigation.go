package telegram

import tele "gopkg.in/telebot.v3"

type NavigationSvc struct {
	MainMenu        *tele.ReplyMarkup
	ChartMenu       *tele.ReplyMarkup
	NotifyMenu      *tele.ReplyMarkup
	BtnChart        tele.Btn
	BtnNotify       tele.Btn
	BtnAddNotify    tele.Btn
	BtnAllNotify    tele.Btn
	BtnCurrentPrice tele.Btn
	BtnChartHour    tele.Btn
	BtnChartDay     tele.Btn
	BtnChartWeek    tele.Btn
	BtnChartMonth   tele.Btn
	BtnPing         tele.Btn
	BtnHome         tele.Btn
}

func NewNavigationSvc() *NavigationSvc {
	mainMenu := &tele.ReplyMarkup{ResizeKeyboard: true}
	chartMenu := &tele.ReplyMarkup{ResizeKeyboard: true}
	notifyMenu := &tele.ReplyMarkup{ResizeKeyboard: true}

	btnNotify := mainMenu.Text("Notify menu")
	btnChart := mainMenu.Text("Charts")
	btnCurrent := mainMenu.Text("Current Price")
	btnPing := mainMenu.Text("Ping Me!")

	btnHour := mainMenu.Text("Hour")
	btnDay := mainMenu.Text("Day")
	btnWeek := mainMenu.Text("Week")
	btnMonth := mainMenu.Text("Month")

	btnHome := mainMenu.Text("Home")
	btnAllNotify := mainMenu.Text("All notifications")
	btnAddNotify := mainMenu.Text("Add notification")

	mainMenu.Reply(
		mainMenu.Row(btnChart, btnNotify),
		mainMenu.Row(btnCurrent, btnPing),
	)

	chartMenu.Reply(
		chartMenu.Row(btnHour, btnDay),
		chartMenu.Row(btnWeek, btnMonth),
		chartMenu.Row(btnHome),
	)

	notifyMenu.Reply(
		notifyMenu.Row(btnAllNotify, btnAddNotify),
		notifyMenu.Row(btnHome),
	)

	return &NavigationSvc{
		MainMenu:        mainMenu,
		ChartMenu:       chartMenu,
		NotifyMenu:      notifyMenu,
		BtnChart:        btnChart,
		BtnNotify:       btnNotify,
		BtnAddNotify:    btnAddNotify,
		BtnAllNotify:    btnAllNotify,
		BtnCurrentPrice: btnCurrent,
		BtnChartHour:    btnHour,
		BtnChartDay:     btnDay,
		BtnChartWeek:    btnWeek,
		BtnChartMonth:   btnMonth,
		BtnPing:         btnPing,
		BtnHome:         btnHome,
	}
}
