package stockhist

import (
	"time"

	mongodb "github.com/williammendat/analytics-hub/MongoDb"
)

type (
	StockHistModel struct {
		mongodb.BaseModel `bson:",inline" json:",inline"`
		Date time.Time `bson:"Date" json:"date"`
		DateString string `bson:"DateString" json:"dateString"`
		Close float64 `bson:"Close" json:"close"`
	}
)

func (s StockHistModel) GetDoc() StockHistModel {
	return s
}