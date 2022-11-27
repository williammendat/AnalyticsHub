package stockrankings

import mongodb "github.com/williammendat/analytics-hub/MongoDb"

type (
	StockPercentModel struct {
		mongodb.BaseModel `bson:",inline" json:",inline"`
		StockPercent `bson:",inline" json:",inline"`
	}
)

func (s StockPercentModel) GetDoc() StockPercentModel {
	return s
}