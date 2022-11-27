package stockinfo

import (
	database "github.com/williammendat/analytics-hub/Database"
	mongodb "github.com/williammendat/analytics-hub/MongoDb"
)

type Repository mongodb.RepositoryI[StockInfoModel, *StockInfoModel]

func NewRepository(dataStore *database.DataStore) Repository {
	return mongodb.NewRepository[StockInfoModel, *StockInfoModel](dataStore.Database.Collection("stockInfos"))
}