package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	log "github.com/sirupsen/logrus"
)

const (
	MongoDbUri      = "MongoDBUri"
	MongoDbDatabase = "MongoDBName"
)

type (
	DataStore struct {
		Client   *mongo.Client
		Database *mongo.Database
		Ctx      context.Context
	}
)

func NewDataStore() (*DataStore, error) {

	mongoDbUri := os.Getenv(MongoDbUri)
	if mongoDbUri == "" {
		msg := fmt.Sprintf("%s is needed as env variable", MongoDbUri)
		return nil, errors.New(msg)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDbUri))
	if err != nil {
		return nil, err
	}

	// Check connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	mongoDbName := os.Getenv(MongoDbDatabase)
	if mongoDbName == "" {
		msg := fmt.Sprintf("%s is needed as env variable", MongoDbDatabase)
		return nil, errors.New(msg)
	}

	db := client.Database(mongoDbName)

	store := &DataStore{	
		Client: client,
		Database: db,
		Ctx: ctx,
	}

	log.Infoln("Connection to MongoDB success")

	return store, nil
}

func (dataStore *DataStore) Disconnect() {

	err := dataStore.Client.Disconnect(dataStore.Ctx)
	if err != nil {
		log.Error(err)
	}

	log.Infoln("Connection to MongoDB closed")
}