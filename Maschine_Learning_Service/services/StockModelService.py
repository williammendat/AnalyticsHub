import numpy as np
import pandas as pd
import xgboost as xgb
from sklearn.metrics import accuracy_score, mean_squared_error
from sklearn.model_selection import train_test_split
from gridfs import Database
import datetime as dt
from dateutil.relativedelta import relativedelta

class StockModelService:

    _client: Database
    _histCollectPrefix: str

    def __init__(self, client: Database) -> None:
        self._client = client
        self._histCollectPrefix = "hist_"
        pass

    def getRelativeStrengthIndex(self, df: pd.DataFrame, n: int = 14):
        close = df['Close']
        delta = close.diff()
        delta = delta.loc[1:]
        pricesUp = delta.copy()
        pricesDown = delta.copy()
        pricesUp[pricesUp < 0] = 0
        pricesDown[pricesDown > 0] = 0
        rollUp = pricesUp.rolling(n).mean()
        rollDown = pricesDown.abs().rolling(n).mean()
        rs = rollUp / rollDown
        rsi = 100.0 - (100.0 / (1.0 + rs))
        return rsi

    def getHistoricalData(self, symbol: str) -> list:
        today = dt.datetime.today()
        start_date = today - relativedelta(years = 3)
        end_date = today - relativedelta(months = 2)

        return self.getHistoricalDataWithDates(symbol, start_date, end_date)

    def getDataframeForPrediction(self, symbol: str, currentPrice) -> pd.DataFrame:
        today = dt.datetime.today()
        start_date = today - relativedelta(months = 4)
        end_date = today + relativedelta(days =  1)

        histList = self.getHistoricalDataWithDates(symbol, start_date, end_date)

        dummyEntry = { "_id": "", "Date": today, "DateString": "", "Close": currentPrice }

        histList.append(dummyEntry)

        return pd.DataFrame(histList) 

    def getHistoricalDataWithDates(self, symbol: str, start_date, end_date) -> list:
        all_collections = self._client.list_collection_names()
        collection = self._histCollectPrefix + symbol
        
        exists = collection in all_collections

        if not exists:
            return list()

        repo = self._client[collection]

        filter = {"Date": {"$gte": start_date, "$lte": end_date}}

        histList = list(repo.find(filter))

        return histList

    def getHistoricalDataAsDataframe(self, symbol: str) -> pd.DataFrame:
        histList = self.getHistoricalData(symbol)

        if len(histList) == 0:
            return pd.DataFrame()

        return pd.DataFrame(histList) 

    def getTechnicalData(self, df: pd.DataFrame) -> pd.DataFrame:
        df['EMA_9'] = df['Close'].ewm(9).mean().shift()
        df['SMA_5'] = df['Close'].rolling(5).mean().shift()
        df['SMA_10'] = df['Close'].rolling(10).mean().shift()
        df['SMA_15'] = df['Close'].rolling(15).mean().shift()
        df['SMA_30'] = df['Close'].rolling(30).mean().shift()

        df.dropna(inplace=True)

        df['RSI'] = self.getRelativeStrengthIndex(df)

        EMA_12 = pd.Series(df['Close'].ewm(span=12, min_periods=12).mean())
        EMA_26 = pd.Series(df['Close'].ewm(span=26, min_periods=26).mean())
        df['MACD'] = pd.Series(EMA_12 - EMA_26)
        df['MACD_signal'] = pd.Series(df["MACD"].ewm(span=9, min_periods=9).mean())

        df.dropna(inplace=True)

        return df

    def addTargetValuesToDataFrame(self, df: pd.DataFrame) -> pd.DataFrame:
        df["Close_5"] = df["Close"].shift(-5)
        df.dropna(inplace=True)
        df["Target"] = np.where(df["Close_5"] > df["Close"], 1, -1)
        return df

    def getDataframeForModel(self, symbol:str) -> pd.DataFrame:
        df = self.getHistoricalDataAsDataframe(symbol)
        
        if len(df) == 0:
            return df

        df = self.getTechnicalData(df)
        df = self.addTargetValuesToDataFrame(df)

        return df

    def getInput(self, df: pd.DataFrame):
        return df[["EMA_9", "SMA_5", "SMA_10", "SMA_15", "SMA_30", "RSI", "MACD", "MACD_signal"]]

    def getInputAndOutput(self, df: pd.DataFrame):
        X = self.getInput(df)
        y = df["Target"]

        return X, y

    def trainAndGetStockModel(self, symbol: str):
        df = self.getDataframeForModel(symbol)

        if len(df) == 0:
            return None, None, None

        X, y = self.getInputAndOutput(df)

        X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=1)

        X_train, X_val, y_train, y_val = train_test_split(X_train, y_train, test_size=0.25, random_state=1)

        best_params = {'gamma': 0.001, 'learning_rate': 0.05, 'max_depth': 12, 'n_estimators': 100, 'random_state': 42}

        eval_set = [(X_train, y_train), (X_val, y_val)]

        model = xgb.XGBRegressor(**best_params, objective='reg:squarederror')
        model.fit(X_train, y_train, eval_set=eval_set, verbose=False)

        y_pred = model.predict(X_test)

        y_pred = np.where(y_pred > 0, 1, -1)

        mse = mean_squared_error(y_test, y_pred)

        accuracy = accuracy_score(y_test, y_pred) * 100.0

        return model, mse, accuracy

    def predict(self, model: xgb.XGBRegressor, symbol: str, currentPrice):
        hist = self.getDataframeForPrediction(symbol, currentPrice)
        hist = self.getTechnicalData(hist)
        hist = hist.tail(1)
        X = self.getInput(hist)
        y = model.predict(X)

        if len(y) == 0:
            return 0

        y = np.where(y > 0, 1, -1)
        return y[0]

