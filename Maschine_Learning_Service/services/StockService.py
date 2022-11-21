from gridfs import Collection, Database
from services.CompanyMarketService import CompanyMarketService
from services.StockModelService import StockModelService
from services.YFinanceService import YFinanceService

import bson
import pandas as pd
from datetime import datetime, timedelta
import pickle


class StockService:
    _histCollectPrefix: str
    _companyMarketService: CompanyMarketService
    _yFinanceService: YFinanceService
    _stockModelService: StockModelService
    _client: Database
    _stocksRepository: Collection
    _stockInfosRepository: Collection
    _syncTaskRepository: Collection
    _syncLogsRepository: Collection
    _stockModelRepository: Collection

    def __init__(self, companyMarketService: CompanyMarketService,
                 yFinanceService: YFinanceService,
                 stockModelService: StockModelService,
                 client: Database,
                 stocksRepository: Collection,
                 stockInfosRepository: Collection,
                 syncTaskRepository: Collection,
                 syncLogsRepository: Collection,
                 stockModelRepository: Collection) -> None:
        self._histCollectPrefix = "hist_"
        self._companyMarketService = companyMarketService
        self._yFinanceService = yFinanceService
        self._stockModelService = stockModelService
        self._client = client
        self._stocksRepository = stocksRepository
        self._stockInfosRepository = stockInfosRepository
        self._syncTaskRepository = syncTaskRepository
        self._syncLogsRepository = syncLogsRepository
        self._stockModelRepository = stockModelRepository
        pass

    def SynchronizeStockInfos(self) -> None:

        syncStockFilter = {"name": "SYNCSTOCKS", "status": 1}

        task = self._syncTaskRepository.find_one(syncStockFilter)

        if task != None:
            print("Sync for stocks is already running")
            return

        syncTaskId = bson.ObjectId()

        syncTask = {
            "_id": syncTaskId,
            "name": "SYNCSTOCKS",
            "status": 1,
            "created_at": datetime.now()
        }

        # Add sync task to db (status 1 (InProgess))
        self._syncTaskRepository.insert_one(syncTask)

        try:
            symbol_dict = self.getSymbolDictionary()

            company_ticker_list = self._companyMarketService.GetCompanyMarketTicker()

            # Create list with new tickers
            ticker_list = list()

            for company_ticker in company_ticker_list:
                if company_ticker["symbol"] in symbol_dict:
                    continue
                ticker_list.append(company_ticker)
                symbol_dict[company_ticker["symbol"]] = 1

            del company_ticker_list
            del symbol_dict

            stocks, stock_infos = self._yFinanceService.getTickerInfos(
                ticker_list)

            del ticker_list

            # save stocks and stock infos
            if len(stocks) > 0:
                self._stocksRepository.insert_many(stocks)

            if len(stock_infos) > 0:
                self._stockInfosRepository.insert_many(stock_infos)

            # sync histories
            self.syncHistory()

            # Status 2 (Finished)
            self.updateTaskAndLog(syncTaskId, 2, f"Success")

        except Exception as e:
            # Status 3 (Failed)
            self.updateTaskAndLog(syncTaskId, 3, f"Exception: {e}")
            print(e)

    def SynchronizeStockModels(self):

        syncModelFilter = {"name": "SYNCSTOCKMODELS", "status": 1}

        task = self._syncTaskRepository.find_one(syncModelFilter)

        if task != None:
            print("Sync for models is already running")
            return

        syncTaskId = bson.ObjectId()

        print("synchronizing models")

        syncTask = {
            "_id": syncTaskId,
            "name": "SYNCSTOCKMODELS",
            "status": 1,
            "created_at": datetime.now()
        }

        # Add sync task to db (status 1 (InProgess))
        self._syncTaskRepository.insert_one(syncTask)

        stocks = list(self._stocksRepository.find({}))

        model_list = list()

        try:
            for stock in stocks:
                print(stock["symbol"])
                symbol = stock["symbol"]

                model, mse, accuracy = self._stockModelService.trainAndGetStockModel(symbol)

                if model == None:
                    print(f'Model for {symbol} failed')
                    continue

                pickle_model = pickle.dumps(model)

                stockModel = {
                    "_id": bson.ObjectId(), 
                    "symbol": symbol,
                    "model": pickle_model, 
                    "mse": mse, 
                    "accuracy": accuracy
                }

                model_list.append(stockModel)

            self._stockModelRepository.delete_many({})
            self._stockModelRepository.insert_many(model_list)

            # Status 2 (Finished)
            self.updateTaskAndLog(syncTaskId, 2, f"Success")
        
        except Exception as e:
            # Status 3 (Failed)
            self.updateTaskAndLog(syncTaskId, 3, f"Exception: {e}")
            print(e)
        

    def updateTaskAndLog(self, taskId: bson.ObjectId, status: int, message: str):
        filter = {"_id": taskId}

        # Update sync task in db
        set = {"$set": {'status': status}}

        self._syncTaskRepository.update_one(filter, set)

        log = {
            "_id": bson.ObjectId(),
            "task_id": taskId,
            "message": message
        }

        self._syncLogsRepository.insert_one(log)

    def getSymbolDictionary(self) -> dict:
        symbolDict = dict()
        stocks = list(self._stocksRepository.find({}))

        for stock in stocks:
            symbolDict[stock["symbol"]] = 1

        return symbolDict

    def syncFullHistory(self, stock):
        symbol = stock["symbol"]
        collection = self._histCollectPrefix + symbol
        hist: pd.DataFrame = self._yFinanceService.GetTickerHistorySubFrame(
            symbol, "max")

        if len(hist) == 0:
            return

        repo = self._client[collection]
        repo.insert_many(hist.to_dict('records'))

    def syncExistingHistory(self, stock):
        symbol = stock["symbol"]
        collection = self._histCollectPrefix + symbol
        repo = self._client[collection]
        last = repo.find_one(sort=[('Date', -1)])

        if last == None:
            return

        date = last["Date"]
        result = date + timedelta(days=1)
        result_string = result.strftime('%Y-%m-%d')

        hist: pd.DataFrame = self._yFinanceService.GetTickerHistoryWithStartDateSubFrame(
            symbol, result_string)
        if len(hist) == 0:
            return

        repo.insert_many(hist.to_dict('records'))

    def syncHistory(self):

        syncHistFilter = {"name": "SYNCHIST", "status": 1}

        task = self._syncTaskRepository.find_one(syncHistFilter)

        if task != None:
            print("Sync for History is already running")
            return

        syncTaskId = bson.ObjectId()

        print("synchronizing models")

        syncTask = {
            "_id": syncTaskId,
            "name": "SYNCHIST",
            "status": 1,
            "created_at": datetime.now()
        }

        # Add sync task to db (status 1 (InProgess))
        self._syncTaskRepository.insert_one(syncTask)

        stocks = list(self._stocksRepository.find({}))
        all_collections = self._client.list_collection_names()

        for stock in stocks:
            print(stock["symbol"])
            symbol = stock["symbol"]
            collection = self._histCollectPrefix + symbol

            exists = collection in all_collections

            if exists:
                self.syncExistingHistory(stock)
            else:
                self.syncFullHistory(stock)

        # Status 2 (Finished)
        self.updateTaskAndLog(syncTaskId, 2, f"Success")

    def getStockData(self, symbol: str) -> dict:
        data = self._yFinanceService.GetFinancialData(symbol)

        model_filter = {"symbol": symbol}

        stock_model = self._stockModelRepository.find_one(model_filter)

        if stock_model == None:
            data["prediction"] = 0
            data["mse"] = 0
            data["accuracy"] = 0

            return data

        model = pickle.loads(stock_model["model"])

        prediction = self._stockModelService.predict(model, symbol, data["currentOpen"])
        
        data["prediction"] = int(prediction)
        data["mse"] = stock_model["mse"]
        data["accuracy"] = stock_model["accuracy"]

        return data
