package stockhist

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	httpclient "github.com/williammendat/analytics-hub/HttpClient"
)

const (
	analyticsHubUrl = "AnalyticsHubUrl"
	stockHistData   = "stockHist/"
)

type (
	HttpClient struct {
		httpClient *httpclient.HttpClient
		baseUrl    string
	}
)

func NewHttpClient(client *httpclient.HttpClient) *HttpClient {
	url := os.Getenv(analyticsHubUrl)
	if url == "" {
		msg := fmt.Sprintf("%s is needed as env variable", analyticsHubUrl)
		logrus.Errorln(errors.New(msg))
	}

	return &HttpClient{
		httpClient: client,
		baseUrl:    url,
	}
}

func (client *HttpClient) GetStockHistorie(symbol, period, interval string) (StockHist, error) {
	url := client.baseUrl + stockHistData + symbol + "/" + period + "/" + interval
	result, err := httpclient.Get[StockHist](client.httpClient, url, map[string]string{})
	return result, err
}