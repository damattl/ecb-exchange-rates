package main

type ExchangeRate struct {
	Currency string `json:"currency" bson:"currency"`
	Rate     string `json:"rate" bson:"rate"`
	Date     int64  `json:"date" json:"date"`
}

type ExchangeRatesForDate struct {
	Date          int64             `json:"date" bson:"date"`
	ExchangeRates map[string]string `json:"exchange_rates" bson:"exchange_rates"`
}
