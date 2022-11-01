package crawler

import (
	"fmt"
	"go.uber.org/zap"
	alert "p2p_crawler/pkg/alert/service"
	"p2p_crawler/pkg/config"
	"p2p_crawler/pkg/price"
	"p2p_crawler/pkg/price/service"
	"time"
)

type Crawler struct {
	exchanges []ExchangeInterface
	priceSvc  *service.PriceService
	alertSvc  *alert.AlertService
	pubSub    PubSubInterface
	logger    *zap.SugaredLogger
	cfg       *config.Config
}

func NewCrawler(
	exch []ExchangeInterface,
	priceSvc *service.PriceService,
	alertSvc *alert.AlertService,
	pubSub PubSubInterface,
	logger *zap.SugaredLogger,
	cfg *config.Config,
) *Crawler {
	return &Crawler{exchanges: exch, priceSvc: priceSvc, alertSvc: alertSvc, pubSub: pubSub, logger: logger, cfg: cfg}
}

type ExchangeInterface interface {
	GetName() string
	DoRequest() (*price.PriceHistory, error)
}

type PubSubInterface interface {
	//Subscribe(topic string) error
	//Unsubscribe(topic string) error
}

func (svc *Crawler) Run() error {
	times := time.Tick(svc.cfg.RequestPeriod)
	for {
		select {
		case <-times:
			for _, exchange := range svc.exchanges {
				t := time.Now()
				svc.logger.Infof("Start requset to %s", exchange.GetName())
				resp, err := exchange.DoRequest()
				if err != nil {
					return fmt.Errorf("exchange request: %w", err)
				}
				tEnd := time.Since(t)

				t1 := time.Now()
				err = svc.priceSvc.Insert(resp)
				if err != nil {
					return fmt.Errorf("save response: %w", err)
				}
				t1End := time.Since(t1)

				t2 := time.Now()
				err = svc.alertSvc.HandlePrice(resp)
				if err != nil {
					return fmt.Errorf("alerting: %w", err)
				}
				t2End := time.Since(t2)

				svc.logger.Infow("Request timings",
					"request", tEnd,
					"insert", t1End,
					"alert", t2End)
			}
		}
	}
}
