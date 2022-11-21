from flask_restful import Resource

from services.StockService import StockService


class StockSync(Resource):

    _stockService: StockService

    def __init__(self, **kwargs):
        # smart_engine is a black box dependency
        self._stockService = kwargs['stockService']
        super().__init__()

    def get(self, sync: str):
        if sync == "full":
            self._stockService.SynchronizeStockInfos()
        elif sync == "hist":
            self._stockService.syncHistory()
        return {"message": "success"}