package main

import (
	"context"
	"errors"
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
	client := getDatabaseClient(appCtx)
	saveExchangeRatesToDB(todaysRates, client) // TODO: Handle error here?

	// TODO: REMOVE LATER
	println(findRateForCurrency("2022-09-21", "HKD", client))

	return nil
}
