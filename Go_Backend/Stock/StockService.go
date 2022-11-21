package stock

import (
	"time"

	"github.com/sirupsen/logrus"
)

type (
	Service struct {
	}
)

func NewService() *Service {
	return &Service{}
}

func (service *Service) SyncStocks() {

	for {
		logrus.Infoln("Sync")

		time.Sleep(1 * time.Second)
	}
}