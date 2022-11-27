package stockrankings

import (
	database "github.com/williammendat/analytics-hub/Database"
	mongodb "github.com/williammendat/analytics-hub/MongoDb"
)

type Repository mongodb.RepositoryI[StockPercentModel, *StockPercentModel]

func NewRepository(dataStore *database.DataStore) Repository {
	return mongodb.NewRepository[StockPercentModel, *StockPercentModel](dataStore.Database.Collection("stockPercent"))
}