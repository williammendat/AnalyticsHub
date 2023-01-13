import pandas as pd
import yfinance as yf
import bson
from deep_translator import GoogleTranslator

class YFinanceService:

    required_ticker_infos = [
        "zip",
        "sector",
        "fullTimeEmployees",
        "longBusinessSummary",
        "city",
        "phone",
        "country",
        "website",
        "address1",
        "industry",
        "regularMarketPrice",
        "logo_url"
    ]

    ticker_info_field = [
        "zip",
        "sector",
        "fullTimeEmployees",
        "city",
        "phone",
        "country",
        "website",
        "industry",
        "logo_url"
    ]

    descriptionKey = "longBusinessSummary"
    address = "address1"

    def __init__(self) -> None:
        pass

    def generateStockObjects(self, ticker_infos):
        stocks = list()
        stock_infos = list()

        for info in ticker_infos:
            stock = {
                "_id": bson.ObjectId(),
                "symbol": info["symbol"],
                "name": info["name"],
                "logo_url": info["logo_url"]
            }

            stock_info = {
                "_id": bson.ObjectId(),
                "ticker_id": stock["_id"],
                "symbol": info["symbol"],
                "name": info["name"],
                "descriptions": info["descriptions"],
                "address": info["address"]
            }

            for info_field in self.ticker_info_field:
                stock_info[info_field] = info[info_field]

            stocks.append(stock)
            stock_infos.append(stock_info)

        return stocks, stock_infos

    def getDescriptions(self, description) -> dict:
        de = GoogleTranslator(
            source='auto', target='de').translate(description)
        fr = GoogleTranslator(
            source='auto', target='fr').translate(description)

        return {"en": description, "de": de, "fr": fr}

    def getTickerInfos(self, ticker_list):
        ticker_infos = list()
        index = 0

        for ticker_object in ticker_list:
            try:
                symbol = ticker_object["symbol"]
                print(symbol)
                tickerInfo = yf.Ticker(symbol)

                info = tickerInfo.info

                has_required_infos = True

                for required_info in self.required_ticker_infos:
                    if (info[required_info] == None or info[required_info] == ""):
                        has_required_infos = False                    

                if not has_required_infos:
                    print(f'{symbol} doesnt have all requirements')
                    continue

                descriptions = self.getDescriptions(info[self.descriptionKey])

                data = {
                    "symbol": ticker_object["symbol"],
                    "name": ticker_object["name"],
                    "descriptions": descriptions,
                    "address": info[self.address]
                }

                for info_field in self.ticker_info_field:
                    data[info_field] = info[info_field]

                ticker_infos.append(data)

            except Exception as e:
                print(e)

            finally:
                index += 1
                print(f'{index}/{len(ticker_list)}')

        return self.generateStockObjects(ticker_infos)

    def GetTickerHistory(self, ticker: str, period: str, inteval: str = "1d") -> pd.DataFrame:
        tickerInfo = yf.Ticker(ticker)
        hist: pd.DataFrame = tickerInfo.history(
            period=period, interval=inteval)
        hist = hist.reset_index()
        hist.dropna(inplace=True)
        return hist

    def GetTickerHistoryWithStartDate(self, ticker: str, start: str, inteval: str = "1d") -> pd.DataFrame:
        tickerInfo = yf.Ticker(ticker)
        hist: pd.DataFrame = tickerInfo.history(start=start, interval=inteval)
        hist = hist.reset_index()
        return hist

    def GetTickerHistoryWithStartDateSubFrame(self, ticker: str, start: str, inteval: str = "1d") -> pd.DataFrame:
        hist = self.GetTickerHistory(ticker, start, inteval)
        return self.GetTickerHistorySubFrameFromHist(hist)

    def GetGetTickerHistoryDateTime(self, ticker: str, period: str, inteval: str = "1d") -> pd.DataFrame:
        hist = self.GetTickerHistory(ticker, period, inteval)
        if "Date" in hist.columns:
            return self.GetTickerHistorySubFrameFromHist(hist)
        return self.GetTickerHistorySubFrameFromHistDateTime(hist)

    def GetTickerHistorySubFrame(self, ticker: str, period: str, inteval: str = "1d") -> pd.DataFrame:
        hist = self.GetTickerHistory(ticker, period, inteval)
        return self.GetTickerHistorySubFrameFromHist(hist)

    def GetTickerHistorySubFrameFromHist(self, hist: pd.DataFrame) -> pd.DataFrame:
        if len(hist) == 0:
            return hist
        hist["DateString"] = hist["Date"].dt.strftime('%Y-%m-%d')
        hist = hist[["Date", "DateString", "Close"]]
        return hist

    def GetTickerHistorySubFrameFromHistDateTime(self, hist: pd.DataFrame) -> pd.DataFrame:
        if len(hist) == 0:
            return hist
        hist["DateString"] = hist["Datetime"].dt.strftime("%m/%d/%Y-%H:%M:%S")
        hist = hist[["Datetime", "DateString", "Close"]]
        return hist

    def GetFinancialData(self, ticker: str) -> dict:
        tickerInfo = yf.Ticker(ticker)
        info = tickerInfo.info
        diff = round(info["currentPrice"] - info["previousClose"], 2)
        percent = round((diff / info["previousClose"]) * 100, 2)

        data = {
            "currentOpen": info["currentPrice"],
            "previousClose": info["previousClose"],
            "diff": diff,
            "diffPercent": percent
        }

        return data
