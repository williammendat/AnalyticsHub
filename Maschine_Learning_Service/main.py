from Container.Container import create_container
from flask import Flask, abort, request
from flask_restful import Api, Resource
from Resources.StockData import StockData
from Resources.StockModel import StockModel
from Resources.StockSync import StockSync
from Resources.YFinanceData import YFinanceData
from Resources.Clear import ClearResource

app = Flask(__name__)
api = Api(app)

trusted_ips = ['127.0.0.1']

@app.before_request
def limit_remote_addr():
    if request.remote_addr not in trusted_ips:
        abort(403) # Forbidden 

container = create_container()

api.add_resource(StockData, "/stockData/<string:symbol>/<string:type>",
                 resource_class_kwargs={'stockService': container["stockService"]})

api.add_resource(YFinanceData, "/stockHist/<string:symbol>/<string:period>/<string:interval>",
                 resource_class_kwargs={'yFinanceService': container["yFinanceService"]})

api.add_resource(StockSync, "/syncStockData/<string:sync>",
                 resource_class_kwargs={'stockService': container["stockService"]})

api.add_resource(StockModel, "/syncModels",
                 resource_class_kwargs={'stockService': container["stockService"]})

api.add_resource(ClearResource, "/clear", 
                resource_class_kwargs={'clearService': container["clearService"]})

if __name__ == "__main__":
    app.run(debug=True)

