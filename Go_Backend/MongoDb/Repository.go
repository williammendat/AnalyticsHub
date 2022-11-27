package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	FindOne[T any] interface {
		FindOne(filter bson.M, opts ...*options.FindOneOptions) (T, error)
	}

	FindMany[T any] interface {
		FindMany(filter bson.M, opts ...*options.FindOptions) ([]T, error)
	}

	InsertOne[T any, D Document[T]] interface {
		InsertOne(doc D, opts ...*options.InsertOneOptions) (T, error)
	}

	InsertMany[T any, D Document[T]] interface {
		InsertMany(docs []D, opts ...*options.InsertManyOptions) ([]T, error)
	}

	UpdateOne interface {
		UpdateOne(filter bson.M, data primitive.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	}

	UpdateMany interface {
		UpdateMany(filter bson.M, data primitive.M, opts ...*options.UpdateOptions)  error
	}

	ReplaceOne[T any, D Document[T]] interface {
		ReplaceOne(filter bson.M, doc D, opts ...*options.ReplaceOptions) (T, error)
	}

	DeleteOne interface {
		DeleteOne(filter bson.M, opts ...*options.DeleteOptions) error
	}

	DeleteMany interface {
		DeleteMany(filter bson.M, opts ...*options.DeleteOptions) (int, error)
	}

	RepositoryI[T any, D Document[T]] interface {
		FindOne[T]
		FindMany[T]
		InsertOne[T, D]
		InsertMany[T, D]
		UpdateOne
		UpdateMany
		ReplaceOne[T, D]
		DeleteOne
		DeleteMany
	}

	Repository[T any, D Document[T]] struct {
		db *mongo.Collection
	}
)

func NewRepository[T any, D Document[T]](collection *mongo.Collection) RepositoryI[T, D] {
	return &Repository[T, D]{
		db: collection,
	}
}

func (r *Repository[T, D]) FindOne(filter bson.M, opts ...*options.FindOneOptions) (T, error) {
	var res T 
	err := r.db.FindOne(context.TODO(), filter, opts...).Decode(&res)

	return res, err
}

func (r *Repository[T, D]) FindMany(filter bson.M, opts ...*options.FindOptions) ([]T, error) {
	var res []T
	cur, err := r.db.Find(context.TODO(), filter, opts...)

	if err != nil {
		return nil, err
	}

	cur.All(context.TODO(), &res)

	return  res, nil
}

func (r *Repository[T, D]) InsertOne(doc D, opts ...*options.InsertOneOptions) (T, error){
	doc.InitMongoID(primitive.NewObjectID())

	_, err := r.db.InsertOne(context.TODO(), doc, opts...)

	return doc.GetDoc(), err
}

func (r *Repository[T, D]) InsertMany(docs []D, opts ...*options.InsertManyOptions) ([]T, error) {
	newDocs := make([]interface{}, len(docs))
	copyDocs := make([]T, len(docs))

	for i, doc := range docs {
		doc.InitMongoID(primitive.NewObjectID())

		newDocs[i] = doc
		copyDocs[i] = doc.GetDoc()
	}

	_, err := r.db.InsertMany(context.TODO(), newDocs, opts...)

	if err != nil {
		return nil, err
	}

	return copyDocs, nil
}

func (r *Repository[T, D]) UpdateOne(filter bson.M, data primitive.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error){
	updateResult, err := r.db.UpdateOne(context.TODO(), filter, bson.M{"$set": data, "$currentDate": bson.M{"updatedAt": true}}, opts...)
	if err != nil {
		return nil, err
	}
	return updateResult, err
}

func (r *Repository[T, D]) UpdateMany(filter bson.M, data primitive.M, opts ...*options.UpdateOptions)(error) {
	_, err := r.db.UpdateMany(context.TODO(), filter, bson.M{"$set": data, "$currentDate": bson.M{"updatedAt": true}}, opts...)
	return err
}

func (r *Repository[T, D]) ReplaceOne(filter bson.M, doc D, opts ...*options.ReplaceOptions) (T, error) {
	_, err := r.db.ReplaceOne(context.TODO(), filter, doc, opts...)
	return doc.GetDoc(), err
}

func (r *Repository[T, D]) DeleteOne(filter bson.M, opts ...*options.DeleteOptions) error {
	_, err := r.db.DeleteOne(context.TODO(), filter, opts...)
	return err
}

func (r *Repository[T, D]) DeleteMany(filter bson.M, opts ...*options.DeleteOptions) (int, error) {
	res, err := r.db.DeleteMany(context.TODO(), filter, opts...)
	return int(res.DeletedCount), err
}