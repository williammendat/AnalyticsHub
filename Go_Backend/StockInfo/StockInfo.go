package stockinfo

type (
	StockGeneralInfo struct {
		Symbol            string `json:"symbol" bson:"symbol"`
		Name              string `json:"name" bson:"name"`
		Address           string `json:"address" bson:"address"`
		Zip               string `json:"zip" bson:"zip"`
		Sector            string `json:"sector" bson:"sector"`
		FullTimeEmployees int    `json:"fullTimeEmployees" bson:"fullTimeEmployees"`
		City              string `json:"city" bson:"city"`
		Phone             string `json:"phone" bson:"phone"`
		Country           string `json:"country" bson:"country"`
		Website           string `json:"website" bson:"website"`
		Industry          string `json:"industry" bson:"industry"`
		LogoURL           string `json:"logoURL" bson:"logo_url"`
	}

	StockInfo struct {
		StockGeneralInfo `bson:",inline" json:",inline"`
		Description string `bson:"description" json:"description"`
	}
)