package crawler

import (
	"fmt"
	"go.uber.org/zap"
	"p2p_crawler/pkg/config"
	"p2p_crawler/pkg/price"
	"p2p_crawler/pkg/price/service"
	"time"
)

type Crawler struct {
	Exchanges    []ExchangeInterface
	PriceService *service.PriceService
	PubSub       PubSubInterface
	Logger       *zap.SugaredLogger
	Cfg          *config.Config
}

func NewCrawler(
	exch []ExchangeInterface,
	priceSvc *service.PriceService,
	pubSub PubSubInterface,
	logger *zap.SugaredLogger,
	cfg *config.Config,
) *Crawler {
	return &Crawler{Exchanges: exch, PriceService: priceSvc, PubSub: pubSub, Logger: logger, Cfg: cfg}
}

type ExchangeInterface interface {
	GetName() string
	DoRequest() (*price.PriceHistoryDB, error)
}

type PubSubInterface interface {
}

func (c *Crawler) Run() error {
	times := time.Tick(c.Cfg.RequestPeriod)
	for {
		select {
		case <-times:
			for _, exchange := range c.Exchanges {
				c.Logger.Infof("start requset to %s", exchange.GetName())
				resp, err := exchange.DoRequest()
				if err != nil {
					return fmt.Errorf("exchange request: %w", err)
				}

				err = c.PriceService.Insert(resp)
				if err != nil {
					return fmt.Errorf("save response: %w", err)
				}

				c.Logger.Infof("end requset to %s", exchange.GetName())
			}
		}
	}
}
