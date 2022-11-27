package stockrankings

type (

	StockRanking struct {
		Ranking int `bson:"ranking" json:"ranking"`
		StockPercent `bson:",inline" json:",inline"`
	}

	StockPercent struct {
		Symbol string `bson:"symbol" json:"symbol"`
		Percent float64 `bson:"percent" json:"percent"`
	}
)