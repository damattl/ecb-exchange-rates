package main

type ExchangeRate struct {
	Currency string `json:"currency" bson:"currency"`
	Rate     string `json:"rate" bson:"rate"`
	//	Date     string `json:"date" json:"date"`
}

type ExchangeRatesForDate struct {
	Date          string            `json:"date" bson:"date"`
	ExchangeRates map[string]string `json:"exchange_rates" bson:"exchange_rates"`
}
