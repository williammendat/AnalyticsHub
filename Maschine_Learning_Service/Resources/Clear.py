from flask_restful import Resource
from services.ClearService import ClearService


class ClearResource(Resource):

    _clearService: ClearService

    def __init__(self, **kwargs):
        # smart_engine is a black box dependency
        self._clearService = kwargs['clearService']
        super().__init__()

    def get(self):
        self._clearService.ClearTasks()
        return {"message": "success"}