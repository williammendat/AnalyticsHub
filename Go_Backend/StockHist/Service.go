package stockhist

import (
	"time"

	database "github.com/williammendat/analytics-hub/Database"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	OneDay     = "1d"
	FiveDay    = "5d"
	OneMonth   = "1mo"
	ThreeMonth = "3mo"
	SixMonth   = "6mo"
	OneYear    = "1y"
	TwoYear    = "2y"
	FiveYear   = "5y"
	TenYear    = "10y"
)

type histFunc func(symbol string) StockHist

type (
	Service struct {
		dataStore   *database.DataStore
		client      *HttpClient
		histFuncMap map[string]histFunc
	}
)

func NewService(
	dataStore *database.DataStore,
	client *HttpClient,
) *Service {
	service := Service{
		dataStore: dataStore,
		client: client,
		histFuncMap: make(map[string]histFunc),
	}

	service.histFuncMap[OneDay] = service.OneDayHist
	service.histFuncMap[FiveDay] = service.FiveDayHist
	service.histFuncMap[OneMonth] = service.OneMonthHist
	service.histFuncMap[ThreeMonth] = service.ThreeMonthHist
	service.histFuncMap[SixMonth] = service.SixMonthHist
	service.histFuncMap[OneYear] = service.OneYearHist
	service.histFuncMap[TwoYear] = service.TwoYearHist
	service.histFuncMap[FiveYear] = service.FiveYearHist
	service.histFuncMap[TenYear] = service.TenYearHist

	return &service
}

func (service *Service) GetHistPeriod(symbol, period string) StockHist {

	histFun, ok := service.histFuncMap[period]
	if !ok {
		panic("Period is not supported")
	} 

	return histFun(symbol)
}

func (service *Service) OneDayHist(symbol string) StockHist {
	stockHist, err := service.client.GetStockHistorie(symbol, OneDay, "15m")
	if err != nil {
		panic(err)
	}
	return stockHist
}

func (service *Service) FiveDayHist(symbol string) StockHist {
	stockHist, err := service.client.GetStockHistorie(symbol, FiveDay, "60m")
	if err != nil {
		panic(err)
	}
	return stockHist
}

func (service *Service) OneMonthHist(symbol string) StockHist {
	stockHist, err := service.client.GetStockHistorie(symbol, OneMonth, "1d")
	if err != nil {
		panic(err)
	}
	return stockHist
}

func (service *Service) ThreeMonthHist(symbol string) StockHist {
	stockHist, err := service.client.GetStockHistorie(symbol, ThreeMonth, "1d")
	if err != nil {
		panic(err)
	}
	return stockHist
}

func (service *Service) SixMonthHist(symbol string) StockHist {
	stockHist, err := service.client.GetStockHistorie(symbol, SixMonth, "5d")
	if err != nil {
		panic(err)
	}
	return stockHist
}

func (service *Service) OneYearHist(symbol string) StockHist {
	return service.getHistFromYear(1, symbol)
}

func (service *Service) TwoYearHist(symbol string) StockHist {
	return service.getHistFromYear(2, symbol)
}

func (service *Service) FiveYearHist(symbol string) StockHist {
	return service.getHistFromYear(5, symbol)
}

func (service *Service) TenYearHist(symbol string) StockHist {
	return service.getHistFromYear(10, symbol)
}

func (service *Service) getHistFromYear(year int, symbol string) StockHist {
	now := time.Now()
	startDate := now.AddDate(-year, 0, 0)

	filter := bson.M{
		"Date": bson.M{
			"$gte": startDate,
		},
	}
	return service.getHist(filter, symbol)
}

func (service *Service) getHist(filter bson.M, symbol string) StockHist {

	repo := CreateRepository(service.dataStore, symbol)

	models, err := repo.FindMany(filter)
	if err != nil {
		panic(err)
	}

	result := StockHist{
		Dates: make([]string, 0, len(models)),
		Values: make([]float64, 0, len(models)),
	}

	for _, model := range models {
		result.Dates = append(result.Dates, model.DateString)
		result.Values = append(result.Values, model.Close)
	}

	return result
}



