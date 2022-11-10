package main

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ilyakaznacheev/cleanenv"
	goose "github.com/lutersergei/p2pcrawler/migration"
	"github.com/lutersergei/p2pcrawler/pkg/alert/handler"
	"github.com/lutersergei/p2pcrawler/pkg/alert/repository"
	alrt "github.com/lutersergei/p2pcrawler/pkg/alert/service"
	"github.com/lutersergei/p2pcrawler/pkg/binance/service"
	chrt "github.com/lutersergei/p2pcrawler/pkg/chart/service"
	"github.com/lutersergei/p2pcrawler/pkg/config"
	"github.com/lutersergei/p2pcrawler/pkg/crawler"
	"github.com/lutersergei/p2pcrawler/pkg/price"
	"github.com/lutersergei/p2pcrawler/pkg/price/repo"
	prc "github.com/lutersergei/p2pcrawler/pkg/price/service"
	"github.com/lutersergei/p2pcrawler/pkg/telegram"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func main() {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()
	defer sugar.Sync()

	// parse config
	cfg := &config.Config{}
	err := cleanenv.ReadConfig("configs/.env", cfg)
	if err != nil {
		sugar.Panic("read config", zap.Error(err))
	}

	// sentry init
	err = sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.SentryDSN,
		TracesSampleRate: 1.0,
		AttachStacktrace: true,
		Debug:            true,
	})
	if err != nil {
		sugar.Fatalf("sentry.Init: %s", err)
	}
	defer func() {
		err := recover()

		if err != nil {
			sentry.CurrentHub().Recover(err)
			sentry.Flush(time.Second * 5)
		}
	}()

	// bot settings
	bot, err := getTgBot(cfg)
	if err != nil {
		sugar.Panic("connect to bot: ", zap.Error(err))
	}
	// todo rewrite logger
	bot.Use(middleware.Logger())

	// DATABASE
	gormDB, err := getDB(cfg)
	if err != nil {
		sugar.Panic("gormDB", zap.Error(err))
	}
	mysqlDB, _ := gormDB.DB()
	// db migrate
	if err := goose.Up(mysqlDB); err != nil {
		sugar.Panic("migration", zap.Error(err))
	}

	httpClient := &http.Client{Timeout: time.Second * 15}

	priceRepo := repo.NewPriceRepo(gormDB)
	alertRepo := repository.NewAlertRepo(gormDB)

	alertTgHandler := handler.NewAlertTgHandler(cfg, sugar, bot)

	binanceSvc := service.NewBinanceService(httpClient)
	priceSvc := prc.NewPriceService(priceRepo, sugar, []price.ExchangeInterface{binanceSvc})
	alertSvc := alrt.NewAlertService(alertTgHandler, alertRepo, sugar)
	chartSvc := chrt.NewChartService(sugar, priceRepo)
	navigationSvc := telegram.NewNavigationSvc()

	tgHandler := telegram.NewTgHandler(navigationSvc, priceSvc, chartSvc, alertSvc)
	tgRouter := telegram.NewTgRouter(tgHandler, bot, navigationSvc)
	tgRouter.ApplyRoutes()

	go bot.Start()

	var pb interface{}
	crw := crawler.NewCrawler([]price.ExchangeInterface{binanceSvc}, priceSvc, alertSvc, pb, sugar, cfg)
	err = crw.Run()
	if err != nil {
		sugar.Panic("crawler svc", zap.Error(err))
	}
}

func getTgBot(cfg *config.Config) (*tele.Bot, error) {
	pref := tele.Settings{
		Token:  cfg.TgBotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	return tele.NewBot(pref)
}

func getDB(cfg *config.Config) (*gorm.DB, error) {
	mysqlDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	mysqlDsn += "&charset=utf8mb4"
	mysqlDsn += "&interpolateParams=true"
	mysqlDsn += "&parseTime=true"
	db, err := gorm.Open(mysql.Open(mysqlDsn), &gorm.Config{})
	return db, err
}
