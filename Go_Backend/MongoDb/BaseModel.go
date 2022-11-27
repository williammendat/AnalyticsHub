package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Document[T any] interface {
		InitMongoID(id primitive.ObjectID)
		GetDoc() T
	}

	BaseModel struct {
		MongoID primitive.ObjectID `json:"Id" bson:"_id"`
	}
)

func (b *BaseModel) InitMongoID(id primitive.ObjectID){
	if b.MongoID.IsZero() {
		b.MongoID = id
	}
}

