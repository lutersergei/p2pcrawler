package service

import (
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/lutersergei/p2pcrawler/pkg/price"
	"go.uber.org/zap"
)

type PriceService struct {
	rep       PriceRepository
	logger    *zap.SugaredLogger
	exchanges []price.ExchangeInterface
}

func NewPriceService(rep PriceRepository, logger *zap.SugaredLogger, exchanges []price.ExchangeInterface) *PriceService {
	return &PriceService{rep: rep, logger: logger, exchanges: exchanges}
}

type PriceRepository interface {
	Insert(history *price.PriceModel) error
}

func (svc *PriceService) Insert(model *price.PriceModel) error {
	return svc.rep.Insert(model)
}

func (svc *PriceService) CurrentPrice() ([]price.CurrentPriceResponse, error) {
	var r []price.CurrentPriceResponse
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

		err := retry.Do(requestFunc, retry.OnRetry(func(n uint, err error) {
			svc.logger.Infof("#%d: %s", n, err)
		}))
		if err != nil {
			return nil, fmt.Errorf("err after retry: %w", err)
		}
		r = append(r, price.CurrentPriceResponse{ExchangeName: exchange.GetName(), BestPrice: resp.BestPrice})
	}

	return r, nil
}
