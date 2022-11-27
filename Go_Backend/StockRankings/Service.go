package stockrankings

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	. "github.com/ahmetb/go-linq/v3"
)

type (
	Service struct {
		repository Repository
	}
)

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (service *Service) GetTop5StockRankings() []StockRanking {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "percent", Value: -1}})
	filter := bson.M{}
	return service.GetNStockRankings(5, filter, findOptions)
}

func (service *Service) GetBottom5StockRankings() []StockRanking {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "percent", Value: 1}})
	filter := bson.M{}
	return service.GetNStockRankings(5, filter, findOptions)
}

func (service *Service) GetFavouritesStockRankings(favourites []string) []StockRanking {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "percent", Value: -1}})
	if len(favourites) == 0 {
		return []StockRanking{}
	}
	filter := bson.M{"symbol": bson.M{"$in": favourites}}
	return service.GetNStockRankings(5, filter, findOptions)
}

func (service *Service) GetNStockRankings(n int64, filter bson.M, findOptions *options.FindOptions) []StockRanking {
	findOptions.SetLimit(n)
	return service.GetStockRankings(filter, findOptions)
}

func (service *Service) GetStockRankings(filter bson.M, findOptions *options.FindOptions) []StockRanking {

	models, err := service.repository.FindMany(filter, findOptions)
	if err != nil {
		panic(err)
	}

	var res []StockRanking
	rank := 0

	From(models).SelectT(func(m StockPercentModel) StockRanking{
		rank += 1
		return StockRanking{
			Ranking: rank,
			StockPercent: m.StockPercent,
		}
	}).ToSlice(&res)

	return res
}
