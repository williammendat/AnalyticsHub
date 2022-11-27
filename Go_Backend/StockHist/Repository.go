package stockhist

import (
	database "github.com/williammendat/analytics-hub/Database"
	mongodb "github.com/williammendat/analytics-hub/MongoDb"
)

const (
	prefix = "hist_"
)

type Repository mongodb.RepositoryI[StockHistModel, *StockHistModel]

func CreateRepository(dataStore *database.DataStore, symbol string) Repository {
	collection := prefix + symbol
	return mongodb.NewRepository[StockHistModel, *StockHistModel](dataStore.Database.Collection(collection))
}