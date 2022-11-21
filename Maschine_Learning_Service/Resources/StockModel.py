from flask_restful import Resource

from services.StockService import StockService


class StockModel(Resource):

    _stockService: StockService

    def __init__(self, **kwargs):
        # smart_engine is a black box dependency
        self._stockService = kwargs['stockService']
        super().__init__()

    def get(self):
        self._stockService.SynchronizeStockModels()
        return {"message": "success"}