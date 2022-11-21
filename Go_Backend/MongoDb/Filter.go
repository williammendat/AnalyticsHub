package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	FilterOption interface {
		Apply(primitive.M)
	}
)

func NewFilter(opt ...FilterOption) primitive.M {
	f := primitive.M{}

	for _, v := range opt {
		v.Apply(f)
	}

	return f
}

type withMongoID primitive.ObjectID

func (w withMongoID) Apply(m primitive.M) {
	m["_id"] = primitive.ObjectID(w)
}

func WithMongoID(id primitive.ObjectID) FilterOption{
	return withMongoID(id)
}