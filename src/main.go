package main

import (
	"context"
	"damattl.de/api/currency/database"
	"damattl.de/api/currency/server"
	"damattl.de/api/currency/tasks"
	"fmt"
	"github.com/go-co-op/gocron"
	"time"
)

func main() {
	appCtx := context.Background()

	database.UseDatabase(func(appCtx context.Context) {
		err := tasks.RefreshTodaysRates(appCtx)
		if err != nil {
			fmt.Printf("there was an error getting the most recent rates: %v\n", err)
		}

		location, err := time.LoadLocation("CET")
		if err != nil {
			panic(err)
		}
		scheduler := gocron.NewScheduler(location)
		scheduler.Every(1).Day().At("16:00;18:00;20:00;23:59").Do(func() {
			err = tasks.RefreshTodaysRates(appCtx)
			if err != nil {
				fmt.Printf("there was an error getting the most recent rates: %v\n", err)
			}
		})

		scheduler.StartAsync()

		server.StartWebSever(appCtx)
	}, appCtx)
}
