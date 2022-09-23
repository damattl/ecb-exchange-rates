package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-co-op/gocron"
	"time"
)

const ecbURL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

func main() {
	appCtx := context.Background()

	useDatabase(func(appCtx context.Context) {
		err := refreshTodaysRates(appCtx)
		if err != nil {
			fmt.Printf("there was an error getting the most recent rates: %v\n", err)
		}

		location, err := time.LoadLocation("CET")
		if err != nil {
			panic(err)
		}
		scheduler := gocron.NewScheduler(location)
		scheduler.Every(1).Day().At("16:00;18:00;20:00;23:59").Do(func() {
			err = refreshTodaysRates(appCtx)
			if err != nil {
				fmt.Printf("there was an error getting the most recent rates: %v\n", err)
			}
		})

		scheduler.StartAsync()

		startWebSever(appCtx)
	}, appCtx)
}

func refreshTodaysRates(appCtx context.Context) error {
	todaysRates := getTodaysRates()
	if todaysRates == nil {
		return errors.New("could not fetch today's rates")
	}

	println(todaysRates.Date)

	client := getDatabaseClient(appCtx)
	saveExchangeRatesToDB(todaysRates, client) // TODO: Handle error here?

	// TODO: REMOVE LATER
	today, _ := time.Parse("2006-01-02", "2022-09-22")
	rate, err := findRateForCurrency("HKD", today.Unix(), client)
	if err != nil {
		println(err.Error())
	}
	println(rate)

	return nil
}
