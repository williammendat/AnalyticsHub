package stock

import (
	stockhist "github.com/williammendat/analytics-hub/StockHist"
	stockinfo "github.com/williammendat/analytics-hub/StockInfo"
	stockprediction "github.com/williammendat/analytics-hub/StockPrediction"
)

type (

	StockRankings struct {
		TopRankings []StockRankingData `bson:"topRankings" json:"topRankings"`
		BottomRankings []StockRankingData `bson:"bottomRankings" json:"bottomRankings"`
	}

	StockRankingData struct{
		Ranking int `bson:"ranking" json:"ranking"`
		Percent float64 `bson:"percent" json:"percent"`
		StockData `json:",inline" bson:",inline"`
		Hist stockhist.StockHist `json:"historie" bson:"historie"`
	}

	StockData struct {
		stockinfo.StockInfo `json:",inline" bson:",inline"`
		StockTempData       `json:",inline" bson:",inline"`
	}

	StockTempData struct {
		StockPriceData                  `json:",inline" bson:",inline"`
		stockprediction.StockPrediction `json:",inline" bson:",inline"`
	}

	StockPriceData struct {
		CurrentOpen   float64 `json:"currentOpen" bson:"currentOpen"`
		PreviousClose float64 `json:"previousClose" bson:"previousClose"`
		Diff          float64 `json:"diff" bson:"diff"`
		DiffPercent       float64 `json:"diffPercent" bson:"diffePercent"`
	}
)
