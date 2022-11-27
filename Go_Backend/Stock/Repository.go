package stock


import (
	database "github.com/williammendat/analytics-hub/Database"
	mongodb "github.com/williammendat/analytics-hub/MongoDb"
)

type Repository mongodb.RepositoryI[StockModel, *StockModel]

func NewRepository(dataStore *database.DataStore) Repository {
	return mongodb.NewRepository[StockModel, *StockModel](dataStore.Database.Collection("stocks"))
}