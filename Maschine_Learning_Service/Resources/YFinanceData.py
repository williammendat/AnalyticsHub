from flask_restful import Resource
from services.YFinanceService import YFinanceService


class YFinanceData(Resource):

    _yFinanceService: YFinanceService

    def __init__(self, **kwargs):
        # smart_engine is a black box dependency
        self._yFinanceService = kwargs['yFinanceService']
        super().__init__()

    def get(self, symbol: str, period: str, interval: str):
        data = self._yFinanceService.GetGetTickerHistoryDateTime(symbol, period, interval)
        dates = data["DateString"].tolist()
        values = data["Close"].tolist()
        return {"dates": dates, "values": values}