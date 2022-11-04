package service

import (
	"github.com/lutersergei/p2pcrawler/pkg/alert"
	"github.com/lutersergei/p2pcrawler/pkg/price"
	"go.uber.org/zap"
)

type AlertService struct {
	handler AlertHandler // todo implement multiple handlers
	rep     AlertRepository
	logger  *zap.SugaredLogger
}

func NewAlertService(h AlertHandler, rep AlertRepository, logger *zap.SugaredLogger) *AlertService {
	return &AlertService{handler: h, rep: rep, logger: logger}
}

type AlertHandler interface {
	GetName() string
	Alert(*alert.AlertDB, *price.PriceModel) error
}

type AlertRepository interface {
	Match(*price.PriceModel) []*alert.AlertDB
	Insert(*alert.AlertDB) error
	Deactivate(*alert.AlertDB)
	GetByUsername(string) ([]*alert.AlertDB, error)
}

func (svc *AlertService) AddAlert(a *alert.AlertDB) error {
	return svc.rep.Insert(a)
}

func (svc *AlertService) GetAlertsByUsername(u string) ([]*alert.AlertDB, error) {
	return svc.rep.GetByUsername(u)
}

func (svc *AlertService) HandlePrice(model *price.PriceModel) error {
	alerts := svc.rep.Match(model)

	for _, a := range alerts {
		err := svc.handler.Alert(a, model)
		if err != nil {
			return err
		}

		svc.rep.Deactivate(a)
	}

	return nil
}
