package stockprediction

import (
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	timeFormat = "2006-01-02"
)

type (
	Service struct {
		repository Repository
		insertMutex sync.Mutex
	}
)

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
		insertMutex: sync.Mutex{},
	}
}

func (service *Service) TryGetCachedPrediction(symbol string) (StockPrediction, bool) {
	date := time.Now()

	dateString := date.Format(timeFormat)

	filter := bson.M{"symbol": symbol, "predictionTimeString": dateString}

	prediction, err := service.repository.FindOne(filter)
	if err != nil {
		return StockPrediction{}, false
	}

	return prediction.StockPrediction, true
}

func (service *Service) InsertPredictionCache(symbol string, prediction StockPrediction) {
	date := time.Now()

	dateString := date.Format(timeFormat)

	model := StockPredictionModel{
		Symbol:               symbol,
		PredictionTime:       date,
		PredictionTimeString: dateString,
		StockPrediction:      prediction,
	}

	service.insertMutex.Lock()
	defer service.insertMutex.Unlock()

	_, exists := service.TryGetCachedPrediction(symbol)
	if exists {
		return
	}

	service.repository.InsertOne(&model)

	filter := bson.M{
		"symbol": symbol,
		"predictionTime": bson.M{
			"$lt": date,
		},
	}

	service.repository.DeleteMany(filter)
}
