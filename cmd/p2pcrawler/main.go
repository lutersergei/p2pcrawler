package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	goose "p2p_crawler/migration"
	"p2p_crawler/pkg/alert"
	"p2p_crawler/pkg/alert/handler"
	"p2p_crawler/pkg/alert/repository"
	service2 "p2p_crawler/pkg/alert/service"
	"p2p_crawler/pkg/binance/service"
	"p2p_crawler/pkg/config"
	"p2p_crawler/pkg/crawler"
	"p2p_crawler/pkg/price/repo"
	price "p2p_crawler/pkg/price/service"
	"strconv"
	"strings"
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
		sugar.Fatalw("read config", zap.Error(err))
	}

	// bot settings
	bot, err := getTgBot(cfg)
	if err != nil {
		sugar.Fatalw("connect to bot: ", zap.Error(err))
	}

	// database
	gormDB, err := getDB(cfg)
	if err != nil {
		sugar.Fatalw("gormDB", zap.Error(err))
	}
	mysqlDB, _ := gormDB.DB()

	// migrate
	if err := goose.Up(mysqlDB); err != nil {
		sugar.Fatalw("migration", zap.Error(err))
	}

	httpClient := &http.Client{Timeout: time.Second * 15}

	binanceSvc := service.NewBinanceService(httpClient)
	priceRepo := repo.NewPriceRepo(gormDB)
	alertRepo := repository.NewAlertRepo(gormDB)
	priceSvc := price.NewPriceService(priceRepo, sugar)

	tgHandler := handler.NewTelegramHandler(cfg, sugar, bot)
	alertSvc := service2.NewAlertService(tgHandler, alertRepo, sugar)

	var pb interface{}
	crw := crawler.NewCrawler([]crawler.ExchangeInterface{binanceSvc}, priceSvc, alertSvc, pb, sugar, cfg)

	addTgHandlers(bot, sugar, alertRepo)
	go bot.Start()

	err = crw.Run()
	if err != nil {
		sugar.Fatalw("crawler error: %w", err)
	}
}

func addTgHandlers(bot *tele.Bot, sugar *zap.SugaredLogger, alertRepo service2.AlertRepository) {
	bot.Handle("/ping", func(c tele.Context) error {
		sugar.Infof("Receive msg from telegram")
		return c.Send("Pong!!")
	})

	bot.Handle("/add", func(c tele.Context) error {
		sugar.Infof("Receive msg from telegram")

		value := strings.TrimPrefix(c.Text(), "/add ")
		pr, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		alertRepo.Insert(&alert.AlertDB{
			Price:     pr,
			Exchange:  "binance",
			Username:  c.Message().Sender.Username,
			Status:    alert.Active,
			CreatedAt: time.Time{},
			MoveType:  alert.MoveUP,
			DealType:  alert.Sell,
		})

		return c.Send(fmt.Sprintf("Add alert for: %v", value))
	})
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
