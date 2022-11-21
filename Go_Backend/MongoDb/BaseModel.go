package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Document[T any] interface {
		InitMongoID(id primitive.ObjectID)
		SetUpdatedAt(updatedAt time.Time)
		SetCreatedAt(createdAt time.Time)
		GetDoc() T
	}

	BaseModel struct {
		MongoID primitive.ObjectID `json:"Id" bson:"_id"`
		CreateAt time.Time `json:"CreatedAt" bson:"createdAt"`
		UpdatedAt time.Time `json:"UpdatedAt" bson:"updatedAt"`
	}
)

func (b *BaseModel) InitMongoID(id primitive.ObjectID){
	if b.MongoID.IsZero() {
		b.MongoID = id
	}
}

func (b *BaseModel) SetUpdatedAt(updatedAt time.Time){
	b.UpdatedAt = updatedAt
}

func (b *BaseModel) SetCreatedAt(CreateAt time.Time){
	b.CreateAt = CreateAt
}
