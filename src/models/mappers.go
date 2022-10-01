package models

import "time"

func ExchangeRatesForDateListToDto(rates []ExchangeRatesForDate) []ExchangeRatesForDateDto {
	ratesDto := make([]ExchangeRatesForDateDto, len(rates))
	for i := range rates {
		date := time.Unix(rates[i].Date, 0).Format("2006-01-02")
		ratesDto[i] = ExchangeRatesForDateDto{
			Date:          date,
			ExchangeRates: rates[i].ExchangeRates,
		}
	}
	return ratesDto
}

func ExchangeRatesForDateToDto(rates *ExchangeRatesForDate) *ExchangeRatesForDateDto {
	date := time.Unix(rates.Date, 0).Format("2006-01-02")
	return &ExchangeRatesForDateDto{
		Date:          date,
		ExchangeRates: rates.ExchangeRates,
	}
}
