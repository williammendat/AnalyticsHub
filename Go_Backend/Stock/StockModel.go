package stock

import mongodb "github.com/williammendat/analytics-hub/MongoDb"

type (
	StockModel struct {
		mongodb.BaseModel `bson:",inline" json:",inline"`
		Symbol            string `bson:"symbol" json:"symbol"`
		Name              string `bson:"name" json:"name"`
	}
)

func (s StockModel) GetDoc() StockModel {
	return s
}