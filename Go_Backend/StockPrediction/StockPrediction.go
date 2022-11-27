package stockprediction

type (
	StockPrediction struct {
		Prediction    int     `bson:"prediction" json:"prediction"`
		Mse           float64 `bson:"mse" json:"mse"`
		Accuracy      float64 `bson:"accuracy" json:"accuracy"`
	}
)