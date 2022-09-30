package tasks

import (
	"context"
	"damattl.de/api/currency/database"
	"damattl.de/api/currency/ecb"
	"damattl.de/api/currency/models"
	"errors"
	"fmt"
	"log"
	"time"
)

func GetTodaysRates() *models.ExchangeRatesForDate {
	ecbResults, err := ecb.GetTodaysRates()
	if err != nil {
		log.Print(err)
		return nil
	}
	var todaysRate models.ExchangeRatesForDate
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

func RefreshTodaysRates(appCtx context.Context) error {
	todaysRates := GetTodaysRates()
	if todaysRates == nil {
		return errors.New("could not fetch today's rates")
	}

	println(todaysRates.Date)

	client := database.GetClient(appCtx)
	database.SaveExchangeRates(todaysRates, client) // TODO: Handle error here?

	/* // TODO: REMOVE LATER
	today, _ := time.Parse("2006-01-02", "2022-09-22")
	rate, err := findRateForCurrency("HKD", today.Unix(), client)
	if err != nil {
		println(err.Error())
	}
	println(rate) */

	return nil
}
