from gridfs import Database
from pymongo import MongoClient

def createDatabase(mongoUri: str, mongoDbName: str) -> Database:
 
   # Provide the mongodb atlas url to connect python to mongodb using pymongo
   CONNECTION_STRING = mongoUri
 
   # Create a connection using MongoClient. You can import MongoClient or use pymongo.MongoClient
   client = MongoClient(CONNECTION_STRING)
 
   # Create the database for our example (we will use the same database throughout the tutorial
   return client[mongoDbName]