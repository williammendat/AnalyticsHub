package stockinfo

import (
	mongodb "github.com/williammendat/analytics-hub/MongoDb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	StockInfoModel struct {
		mongodb.BaseModel `bson:",inline" json:",inline"`
		TickerID primitive.ObjectID `bson:"ticker_id" json:"tickerId"`
		Descriptions map[string]string `bson:"descriptions" json:"descriptions"`
		StockGeneralInfo `bson:",inline" json:",inline"`
	}
)

func (s StockInfoModel) GetDoc() StockInfoModel {
	return s
}