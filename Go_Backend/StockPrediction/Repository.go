package stockprediction

import (
	database "github.com/williammendat/analytics-hub/Database"
	mongodb "github.com/williammendat/analytics-hub/MongoDb"
)

type Repository mongodb.RepositoryI[StockPredictionModel, *StockPredictionModel]

func NewRepository(dataStore *database.DataStore) Repository {
	return mongodb.NewRepository[StockPredictionModel, *StockPredictionModel](dataStore.Database.Collection("stockTempPredictions"))
}