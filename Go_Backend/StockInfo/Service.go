package stockinfo

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	Service struct {
		repository Repository
	}
)

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (service *Service) GetStockInfo(symbol, language string) StockInfo {

	filter := bson.M{
		"symbol": symbol,
	}

	model, err := service.repository.FindOne(filter)
	if err != nil {
		logrus.Errorln(err)
		panic(err)
	}

	description := ""

	existingDescription, ok := model.Descriptions[language]
	if ok {
		description = existingDescription
	}

	info := StockInfo{
		StockGeneralInfo: model.StockGeneralInfo,
		Description: description,
	}

	return info
}