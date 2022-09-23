package main

import (
	"context"
	"errors"
	"time"
)

const ecbURL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

func main() {
	appCtx := context.Background()

	useDatabase(func(appCtx context.Context) {
		refreshTodaysRates(appCtx) // TODO: Handle errors

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
