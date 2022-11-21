from gridfs import Collection, Database

class ClearService:
    _syncTaskRepository: Collection
    _syncLogsRepository: Collection

    def __init__(self,
                syncTaskRepository: Collection,
                 syncLogsRepository: Collection,) -> None:
        self._syncTaskRepository = syncTaskRepository
        self._syncLogsRepository = syncLogsRepository
        pass

    def ClearTasks(self):
        filter = {"status": {"$ne" : 1} }
        self._syncTaskRepository.delete_many(filter)
        self._syncLogsRepository.delete_many({})