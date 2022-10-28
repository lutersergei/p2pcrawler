package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	goose "p2p_crawler/migration"
	"p2p_crawler/pkg/alert/handler"
	"p2p_crawler/pkg/alert/repository"
	service2 "p2p_crawler/pkg/alert/service"
	"p2p_crawler/pkg/binance/service"
	"p2p_crawler/pkg/config"
	"p2p_crawler/pkg/crawler"
	"p2p_crawler/pkg/price/repo"
	price "p2p_crawler/pkg/price/service"
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

	// database
	gormDB, err := getMySQL(cfg)
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

	tgHandler := handler.NewTelegramHandler(cfg, sugar)
	alertSvc := service2.NewAlertService(tgHandler, alertRepo, sugar)

	var pb interface{}
	crw := crawler.NewCrawler([]crawler.ExchangeInterface{binanceSvc}, priceSvc, alertSvc, pb, sugar, cfg)

	err = crw.Run()
	if err != nil {
		sugar.Fatalw("crawler error: %w", err)
	}
}

func getMySQL(cfg *config.Config) (*gorm.DB, error) {
	//mysqlDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?", "root", "E94XuDv35Zp6", "p2p_db", "3306", "p2p")
	mysqlDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	mysqlDsn += "&charset=utf8mb4"
	mysqlDsn += "&interpolateParams=true"
	mysqlDsn += "&parseTime=true"
	db, err := gorm.Open(mysql.Open(mysqlDsn), &gorm.Config{})
	return db, err
}
