package stock

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	httpclient "github.com/williammendat/analytics-hub/HttpClient"
)

const (
	analyticsHubUrl = "AnalyticsHubUrl"
	stockDataPath   = "stockData/"
	clearPath       = "clear/"
	stockHistData   = "stockHist"
	syncStockData   = "syncStockData/"
	syncModels      = "syncModels"
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

func (client *HttpClient) GetFullStockData(symbol string) (StockTempData, error) {
	url := client.baseUrl + stockDataPath + symbol + "/full"
	result, err := httpclient.Get[StockTempData](client.httpClient, url, map[string]string{})
	return result, err
}

func (client *HttpClient) GetCurrentStockData(symbol string) (StockPriceData, error) {
	url := client.baseUrl + stockDataPath + symbol + "/price"
	result, err := httpclient.Get[StockPriceData](client.httpClient, url, map[string]string{})
	return result, err
}

func (client *HttpClient) SyncStock(syncType string) (map[string]string, error) {
	url := client.baseUrl + syncStockData + syncType
	result, err := httpclient.Get[map[string]string](client.httpClient, url, map[string]string{})
	return result, err
}

func (client *HttpClient) SyncStockModels() (map[string]string, error) {
	url := client.baseUrl + syncModels 
	result, err := httpclient.Get[map[string]string](client.httpClient, url, map[string]string{})
	return result, err
}

func (client *HttpClient) SyncStockHistorie() (map[string]string, error) {
	return client.SyncStock("hist")
}

func (client *HttpClient) SyncStockFull() (map[string]string, error) {
	return client.SyncStock("full")
}

func (client *HttpClient) ClearSyncTask() (map[string]string, error) {
	url := client.baseUrl + clearPath
	result, err := httpclient.Get[map[string]string](client.httpClient, url, map[string]string{})
	return result, err
}
