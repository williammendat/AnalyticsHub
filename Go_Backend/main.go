package main

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	database "github.com/williammendat/analytics-hub/Database"
	scheduler "github.com/williammendat/analytics-hub/Scheduler"
	stock "github.com/williammendat/analytics-hub/Stock"
	_ "github.com/williammendat/analytics-hub/User"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			logrus.Errorln(err)
		}
	}()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Could not load env")
		fmt.Println(err)
	}

	dataStore, err := database.NewDataStore()
	if err != nil {
		fmt.Println("Could not load env")
		fmt.Println(err)
	}
	defer dataStore.Disconnect()

	syncService := stock.NewService()
	scheduler := scheduler.NewScheduler(syncService.SyncStocks)

	go scheduler.StartTasksAsync()
	
	time.Sleep(10 * time.Second)
}
