package main

type apiErrorCode int32

const CURRENCY_NOT_FOUND_ERROR apiErrorCode = 0
const NO_ENTRY_FOUND_ERROR apiErrorCode = 1
const FUTURE_DATE_ERROR apiErrorCode = 2
const FORMAT_NOT_SUPPORTED_ERROR apiErrorCode = 3

type ExchangeRate struct {
	Currency string `json:"currency" bson:"currency"`
	Rate     string `json:"rate" bson:"rate"`
	Date     int64  `json:"date" json:"date"`
}

type ExchangeRatesForDate struct {
	Date          int64             `json:"date" bson:"date"`
	ExchangeRates map[string]string `json:"exchange_rates" bson:"exchange_rates"`
}

type APIError struct {
	Message string       `json:"message"`
	Code    apiErrorCode `json:"code"`
}
