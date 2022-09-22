package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v\n", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v\n", err)
	}

	return data, nil
}

func getTodaysRates() *ExchangeRatesForDate {

	if xmlBytes, err := getXML(ecbURL); err != nil {
		log.Printf("Failed to get XML: %v\n", err)
	} else {
		var result ECBRatesEnvelope
		if err = xml.Unmarshal(xmlBytes, &result); err != nil {
			log.Fatalf("Could not parse data: %v\n", err)
		}
		var todaysRate ExchangeRatesForDate
		todaysRate.ExchangeRates = make(map[string]string)
		for _, exRate := range result.ExchangeRates.ExchangeRatesForTime.ExchangeRateInfo {
			fmt.Printf("Rate for currency %s is %s today \n", exRate.Currency, exRate.Rate)
			todaysRate.ExchangeRates[exRate.Currency] = exRate.Rate
		}
		todaysRate.Date = result.ExchangeRates.ExchangeRatesForTime.Time
		return &todaysRate
	}
	return nil
}
