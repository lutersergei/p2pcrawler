package service

import (
	"go.uber.org/zap"
	"p2p_crawler/pkg/price"
)

type PriceService struct {
	rep    PriceRepository
	logger *zap.SugaredLogger
}

func NewPriceService(rep PriceRepository, logger *zap.SugaredLogger) *PriceService {
	return &PriceService{rep: rep, logger: logger}
}

type PriceRepository interface {
	Insert(history *price.PriceHistory) error
}

func (svc PriceService) Insert(model *price.PriceHistory) error {
	return svc.rep.Insert(model)
}
