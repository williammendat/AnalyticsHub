package stockhist

type(
	StockHist struct {
		Dates  []string  `json:"dates" bson:"dates"`
		Values []float64 `json:"values" bson:"values"`
	}
)