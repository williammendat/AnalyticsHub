import os

from dotenv import load_dotenv

from Datastore.Database import createDatabase
from services.ClearService import ClearService
from services.CompanyMarketService import CompanyMarketService
from services.StockModelService import StockModelService
from services.StockService import StockService
from services.YFinanceService import YFinanceService


def create_container() -> dict:
    load_dotenv()

    MongoDbUri = os.getenv("MongoDBUri")
    MongoDbName = os.getenv("MongoDBName")
  
    database = createDatabase(MongoDbUri, MongoDbName)
    stocksRepository = database["stocks"]
    stockInfosRepository = database["stockInfos"]
    stockPercentRepository = database["stockPercent"]
    syncTasksRepositoy = database["syncTasks"]
    syncLogsRepository = database["syncLogs"]
    stockModelRepository = database["stockModels"]

    yFinanceService = YFinanceService()
    companyMarketService = CompanyMarketService()
    stockModelService = StockModelService(database)
    clearService = ClearService(syncTasksRepositoy, syncLogsRepository)

    stockService = StockService(companyMarketService, yFinanceService, stockModelService, database,
                                stocksRepository, stockInfosRepository, stockPercentRepository, syncTasksRepositoy, syncLogsRepository, stockModelRepository)

    container = {
        "database": database,
        "stocksRepository": stocksRepository,
        "stockInfosRepository": stockInfosRepository,
        "syncTasksRepositoy": syncTasksRepositoy,
        "syncLogsRepository": syncLogsRepository,
        "stockModelRepository": stockModelRepository,
        "yFinanceService": yFinanceService,
        "companyMarketService": companyMarketService,
        "stockModelService": stockModelService,
        "stockService": stockService,
        "clearService": clearService
    }

    return container