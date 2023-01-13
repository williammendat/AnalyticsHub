package stock

import mongodb "github.com/williammendat/analytics-hub/MongoDb"

type (
	StockModel struct {
		mongodb.BaseModel `bson:",inline" json:",inline"`
		Symbol            string `bson:"symbol" json:"symbol"`
		Name              string `bson:"name" json:"name"`
		LogoURL           string `json:"logoURL" bson:"logo_url"`
	}

	PagedStockModels struct{
		StockModels []StockModel `bson:"models" json:"models"`
		Count int `bson:"count" json:"count"`
		Page int `bson:"page" json:"page"`
		MaxPage int `bson:"maxPage" json:"maxPage"`
	}
)

func (s StockModel) GetDoc() StockModel {
	return s
}