package models

type ExchangeRate struct {
	Currency string  `json:"currency" bson:"currency"`
	Rate     float64 `json:"rate" bson:"rate"`
	Date     int64   `json:"date" json:"date"`
}

type ExchangeRatesForDate struct {
	Date          int64             `json:"date" bson:"date"`
	ExchangeRates map[string]string `json:"exchange_rates" bson:"exchange_rates"` // [currency]rate
}

type ExchangeRateDto struct {
	Currency string  `json:"currency" bson:"currency"`
	Rate     float64 `json:"rate" bson:"rate"`
	Date     string  `json:"date" json:"date"`
}

type ExchangeRatesForDateDto struct {
	Date          string            `json:"date" bson:"date"`
	ExchangeRates map[string]string `json:"exchange_rates" bson:"exchange_rates"` // [currency]rate
}
