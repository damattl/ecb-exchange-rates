package main

import (
	"context"
)

const ecbURL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

func main() {
	appCtx := context.Background()

	useDatabase(func(appCtx context.Context) {
		todaysRates := getTodaysRates()
		if todaysRates == nil {
			// TODO: HANDLE ERROR
		}
		client := getDatabaseClient(appCtx)
		saveExchangeRatesToDB(todaysRates, client)

		println(findRateForCurrency("2022-09-21", "HKD", client))

		startWebSever(appCtx)
	}, appCtx)
}
