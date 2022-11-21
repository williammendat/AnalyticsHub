package user

import (
	database "github.com/williammendat/analytics-hub/Database"
	mongodb "github.com/williammendat/analytics-hub/MongoDb"
)

type Repository mongodb.RepositoryI[UserModel, *UserModel]

func NewRepository(dataStore *database.DataStore) Repository {
	return mongodb.NewRepository[UserModel, *UserModel](dataStore.Database.Collection("users"))
}


