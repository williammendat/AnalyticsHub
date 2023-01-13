package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	auth "github.com/williammendat/analytics-hub/Auth"
	database "github.com/williammendat/analytics-hub/Database"
	httpclient "github.com/williammendat/analytics-hub/HttpClient"
	generalrouter "github.com/williammendat/analytics-hub/Router"
	scheduler "github.com/williammendat/analytics-hub/Scheduler"
	stock "github.com/williammendat/analytics-hub/Stock"
	stockhist "github.com/williammendat/analytics-hub/StockHist"
	stockinfo "github.com/williammendat/analytics-hub/StockInfo"
	stockprediction "github.com/williammendat/analytics-hub/StockPrediction"
	stockrankings "github.com/williammendat/analytics-hub/StockRankings"
	user "github.com/williammendat/analytics-hub/User"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			logrus.Errorln(err)
		}
	}()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Could not load env")
		fmt.Println(err)
	}

	dataStore, err := database.NewDataStore()
	if err != nil {
		fmt.Println("Could not load env")
		fmt.Println(err)
	}
	defer dataStore.Disconnect()

	baseClient := httpclient.NewClient()

	// User
	userRepository := user.NewRepository(dataStore)
	userService := user.NewService(userRepository)

	// Auth
	_ = auth.NewService(userService)

	// Stock Ranking 
	stockPercentRepository := stockrankings.NewRepository(dataStore)
	stockRankingService := stockrankings.NewService(stockPercentRepository)

	// Stock Prediction
	stockPredictionRepository := stockprediction.NewRepository(dataStore)
	stockPredictionService := stockprediction.NewService(stockPredictionRepository)

	// Stock Info
	stockInfoRepository := stockinfo.NewRepository(dataStore)
	stockInfoService := stockinfo.NewService(stockInfoRepository)

	// Stock Hist
	stockHistHttpClient := stockhist.NewHttpClient(baseClient)
	stockHistService := stockhist.NewService(dataStore, stockHistHttpClient)
	stockHistController := stockhist.NewController(stockHistService)
	stockHistRouter := stockhist.NewRouter(stockHistController)

	// Stock
	stockHttpClient := stock.NewHttpClient(baseClient)
	stockRepository := stock.NewRepository(dataStore)
	stockService := stock.NewService(stockHttpClient, stockRepository, stockPredictionService, stockInfoService, stockRankingService, stockHistService)
	stockController := stock.NewController(stockService)
	stockRouter := stock.NewRouter(stockController)


	// Scheduler
	scheduler := scheduler.NewScheduler(
		stockService.SyncStocks,
		stockService.SyncStockHists,
		stockService.SyncStockModels,
		stockService.ClearSyncTask,
	)

	go scheduler.StartTasksAsync()

	router := gin.Default()
	stockHistRouter.InitRouter(router)
	stockRouter.InitRouter(router)
    
	generalrouter.InitRouter(router)

    router.Run("localhost:2345")
}
