from flask_restful import Resource
import json 

from services.StockService import StockService


class StockData(Resource):

    _stockService: StockService

    def __init__(self, **kwargs):
        # smart_engine is a black box dependency
        self._stockService = kwargs['stockService']
        super().__init__()

    def get(self, symbol: str, type: str):
        withPrediction = False
        if type == "full":
            withPrediction = True
        data = self._stockService.getStockData(symbol, withPrediction)
        return json.dumps(data)