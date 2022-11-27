package stock

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	stockhist "github.com/williammendat/analytics-hub/StockHist"
	stockinfo "github.com/williammendat/analytics-hub/StockInfo"
	stockprediction "github.com/williammendat/analytics-hub/StockPrediction"
	stockrankings "github.com/williammendat/analytics-hub/StockRankings"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	Service struct {
		stockClient            *HttpClient
		repository             Repository
		stockPredictionService *stockprediction.Service
		stockInfoService       *stockinfo.Service
		stockRankingService    *stockrankings.Service
		stockHistService       *stockhist.Service
	}
)

func NewService(
	stockClient *HttpClient,
	repository Repository,
	stockPredictionService *stockprediction.Service,
	stockInfoService *stockinfo.Service,
	stockRankingService *stockrankings.Service,
	stockHistService *stockhist.Service,
) *Service {
	return &Service{
		stockClient:            stockClient,
		repository:             repository,
		stockPredictionService: stockPredictionService,
		stockInfoService:       stockInfoService,
		stockRankingService:    stockRankingService,
		stockHistService:       stockHistService,
	}
}

func (service *Service) GetStockModels() []StockModel {
	filter := bson.M{}

	models, err := service.repository.FindMany(filter)
	if err != nil {
		panic(err)
	}

	return models
}

func (service *Service) GetStockRankings(withHist bool) StockRankings {
	stockRankings := StockRankings{
		TopRankings:    make([]StockRankingData, 5),
		BottomRankings: make([]StockRankingData, 5),
	}

	wgRankings := sync.WaitGroup{}

	var tops []stockrankings.StockRanking
	wgRankings.Add(1)
	go func() {
		tops = service.stockRankingService.GetTop5StockRankings()
		wgRankings.Done()
	}()

	var bottoms []stockrankings.StockRanking
	wgRankings.Add(1)
	go func() {
		bottoms = service.stockRankingService.GetBottom5StockRankings()
		wgRankings.Done()
	}()

	wgRankings.Wait()

	wgTops := sync.WaitGroup{}

	if len(tops) != 0 {
		for i, top := range tops {
			wgTops.Add(1)
			go service.addStockRankingData(i, &wgTops, top, &stockRankings.TopRankings)
		}
	}

	wgBottoms := sync.WaitGroup{}

	if len(bottoms) != 0 {
		for i, bottom := range bottoms {
			wgBottoms.Add(1)
			go service.addStockRankingData(i, &wgBottoms, bottom, &stockRankings.BottomRankings)
		}
	}

	wgTops.Wait()
	wgBottoms.Wait()

	if !withHist {
		return stockRankings
	}

	wgTopHist := sync.WaitGroup{}

	for i, topRanking := range stockRankings.TopRankings {
		wgTopHist.Add(1)
		go func(index int, stockData StockRankingData){
			stockRankings.TopRankings[index].Hist = service.stockHistService.GetHistPeriod(stockData.Symbol, stockhist.FiveDay)
			wgTopHist.Done()
		}(i, topRanking)
	}

	wgBottomHist := sync.WaitGroup{}

	for i, bottomRanking := range stockRankings.BottomRankings {
		wgBottomHist.Add(1)
		go func(index int, stockData StockRankingData){
			stockRankings.BottomRankings[index].Hist = service.stockHistService.GetHistPeriod(stockData.Symbol, stockhist.FiveDay)
			wgBottomHist.Done()
		}(i, bottomRanking)
	}

	wgTopHist.Wait()
	wgBottomHist.Wait()

	return stockRankings
}


func (service *Service) addStockRankingData(index int, wg *sync.WaitGroup, ranking stockrankings.StockRanking, rankings *[]StockRankingData) {
	stockData := service.GetStockData(ranking.Symbol)

	stockRankingData := StockRankingData{
		Ranking:   ranking.Ranking,
		Percent:   ranking.Percent,
		StockData: stockData,
		Hist: stockhist.StockHist{
			Dates: make([]string, 0),
			Values: make([]float64, 0),
		},
	}

	(*rankings)[index] = stockRankingData

	wg.Done()
}

func (service *Service) GetStockData(symbol string) StockData {
	stockTempData := StockTempData{}
	stockInfo := stockinfo.StockInfo{}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		stockTempData = service.GetStockTempData(symbol)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		stockInfo = service.stockInfoService.GetStockInfo(symbol, "de")
		wg.Done()
	}()

	wg.Wait()

	return StockData{
		StockInfo:     stockInfo,
		StockTempData: stockTempData,
	}
}

func (service *Service) GetStockTempData(symbol string) StockTempData {
	stockTempData := StockTempData{}

	stockPrediction, ok := service.stockPredictionService.TryGetCachedPrediction(symbol)
	if ok {
		stockPriceData, err := service.stockClient.GetCurrentStockData(symbol)
		if err != nil {
			panic(err)
		}
		stockTempData = StockTempData{
			StockPriceData:  stockPriceData,
			StockPrediction: stockPrediction,
		}
	} else {
		tempData, err := service.stockClient.GetFullStockData(symbol)
		if err != nil {
			panic(err)
		}
		go service.stockPredictionService.InsertPredictionCache(symbol, tempData.StockPrediction)
		stockTempData = tempData
	}

	return stockTempData
}

func (service *Service) SyncStocks() {
	for {
		now := time.Now()
		nextSync := now.AddDate(0, 1, 0)

		logrus.Infoln("Next full sync for stocks will be on:" + nextSync.String())

		time.Sleep(time.Until(nextSync))
		_, err := service.stockClient.SyncStockFull()
		if err != nil {
			logrus.Infoln("Full sync failed:" + err.Error())
		}
	}
}

func (service *Service) SyncStockHists() {
	for {
		logrus.Infoln("History sync started")
		_, err := service.stockClient.SyncStockHistorie()
		if err != nil {
			logrus.Errorln("Hist sync failed:" + err.Error())
		}
		logrus.Infoln("History sync finisched")

		now := time.Now()
		nextSync := now.AddDate(0, 0, 1)

		logrus.Infoln("Next history sync for stocks will be on:" + nextSync.String())

		time.Sleep(time.Until(nextSync))
	}
}

func (service *Service) SyncStockModels() {
	for {
		now := time.Now()
		nextSync := now.AddDate(0, 2, 1)

		logrus.Infoln("Next sync for stocks model will be on:" + nextSync.String())

		time.Sleep(time.Until(nextSync))

		logrus.Infoln("Model sync started")
		_, err := service.stockClient.SyncStockModels()
		if err != nil {
			logrus.Errorln("Model sync failed:" + err.Error())
		}
		logrus.Infoln("Model sync finisched")
	}
}

func (service *Service) ClearSyncTask() {
	for {
		now := time.Now()
		nextSync := now.AddDate(0, 3, 3)

		logrus.Infoln("Next clear will be on:" + nextSync.String())

		time.Sleep(time.Until(nextSync))

		logrus.Infoln("clear started")
		_, err := service.stockClient.ClearSyncTask()
		if err != nil {
			logrus.Errorln("clear failed:" + err.Error())
		}
		logrus.Infoln("clear finisched")
	}
}

