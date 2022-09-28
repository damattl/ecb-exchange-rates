package main

import (
	"context"
	"damattl.de/api/currency/ecb"
	"errors"
	"fmt"
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

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

func getTodaysRates() *ExchangeRatesForDate {
	ecbResults, err := ecb.GetTodaysRates()
	if err != nil {
		log.Print(err)
		return nil
	}
	var todaysRate ExchangeRatesForDate
	todaysRate.ExchangeRates = make(map[string]string)
	for _, exRate := range ecbResults.ExchangeRates.ExchangeRatesForTime.ExchangeRateInfo {
		fmt.Printf("Rate for currency %s is %s today \n", exRate.Currency, exRate.Rate)
		todaysRate.ExchangeRates[exRate.Currency] = exRate.Rate
	}
	date, err := time.Parse("2006-01-02", ecbResults.ExchangeRates.ExchangeRatesForTime.Time)
	if err != nil {
		log.Fatalf("Could not parse date: %v", err)
	}
	todaysRate.Date = date.Unix()
	return &todaysRate
}

func refreshTodaysRates(appCtx context.Context) error {
	todaysRates := getTodaysRates()
	if todaysRates == nil {
		return errors.New("could not fetch today's rates")
	}

	println(todaysRates.Date)

	client := getDatabaseClient(appCtx)
	saveExchangeRatesToDB(todaysRates, client) // TODO: Handle error here?

	/* // TODO: REMOVE LATER
	today, _ := time.Parse("2006-01-02", "2022-09-22")
	rate, err := findRateForCurrency("HKD", today.Unix(), client)
	if err != nil {
		println(err.Error())
	}
	println(rate) */

	return nil
}
