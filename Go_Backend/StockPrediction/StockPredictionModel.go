package stockprediction

import (
	"time"

	mongodb "github.com/williammendat/analytics-hub/MongoDb"
)

type (
	StockPredictionModel struct {
		mongodb.BaseModel `bson:",inline" json:",inline"`
		Symbol string `bson:"symbol" json:"symbol"`
		PredictionTime time.Time `bson:"predictionTime" json:"predictionTime"`
		PredictionTimeString string `bson:"predictionTimeString" json:"predictionTimeString"`
		StockPrediction `bson:",inline" json:",inline"`
	}
)

func (s StockPredictionModel) GetDoc() StockPredictionModel {
	return s
}