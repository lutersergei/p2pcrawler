package crawler

import (
	"fmt"
	"github.com/avast/retry-go/v4"
	"go.uber.org/zap"
	alert "p2p_crawler/pkg/alert/service"
	"p2p_crawler/pkg/config"
	"p2p_crawler/pkg/price"
	"p2p_crawler/pkg/price/service"
	"time"
)

type Crawler struct {
	exchanges []price.ExchangeInterface
	priceSvc  *service.PriceService
	alertSvc  *alert.AlertService
	pubSub    price.PubSubInterface
	logger    *zap.SugaredLogger
	cfg       *config.Config
}

func NewCrawler(
	exch []price.ExchangeInterface,
	priceSvc *service.PriceService,
	alertSvc *alert.AlertService,
	pubSub price.PubSubInterface,
	logger *zap.SugaredLogger,
	cfg *config.Config,
) *Crawler {
	return &Crawler{exchanges: exch, priceSvc: priceSvc, alertSvc: alertSvc, pubSub: pubSub, logger: logger, cfg: cfg}
}

func (svc *Crawler) Run() error {
	times := time.Tick(svc.cfg.RequestPeriod)
	for {
		select {
		case <-times:
			for _, exchange := range svc.exchanges {
				var resp *price.PriceModel

				var requestFunc retry.RetryableFunc = func() error {
					var err error
					resp, err = exchange.DoRequest()
					if err != nil {
						return fmt.Errorf("exchange request: %w", err)
					}
					return nil
				}

				t := time.Now()

				svc.logger.Infof("Start request to %s", exchange.GetName())
				err := retry.Do(requestFunc, retry.OnRetry(func(n uint, err error) {
					svc.logger.Infof("#%d: %s", n, err)
				}))
				if err != nil {
					return fmt.Errorf("err after retry: %w", err)
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
