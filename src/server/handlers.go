package server

import (
	"context"
	"damattl.de/api/currency/database"
	"damattl.de/api/currency/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strconv"
	"time"
)

func getSupportedCurrencies(w http.ResponseWriter, r *http.Request, appCtx context.Context) {
	client, ok := appCtx.Value(database.MONGO_DB_CLIENT).(*mongo.Client)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Database Client not found")
		return
	}

	unixDate := int64(-1)
	urlVars := mux.Vars(r)
	date, ok := urlVars["date"]
	if ok {
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeError(models.APIError{"date format not supported", models.FORMAT_NOT_SUPPORTED_ERROR}, w)
			return
		}
		today := time.Now()
		isFuture := today.Before(parsedDate)
		if isFuture {
			w.WriteHeader(http.StatusBadRequest)
			writeError(models.APIError{"date is in the future", models.FUTURE_DATE_ERROR}, w)
			return
		}
		unixDate = parsedDate.Unix()
	}

	supportedCurrencies, err := database.FindAllSupportedCurrencies(unixDate, client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			writeError(models.APIError{"could not find any information", models.NO_ENTRY_FOUND_ERROR}, w)
			return
		}
		log.Println(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(supportedCurrencies)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getRatesUntil(w http.ResponseWriter, r *http.Request, appCtx context.Context) {
	client, ok := appCtx.Value(database.MONGO_DB_CLIENT).(*mongo.Client)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Database Client not found")
		return
	}

	urlVars := mux.Vars(r)

	unixDate, isFuture, err := parseDateAndHandleError(w, r, urlVars)
	if err != nil {
		return
	}

	if isFuture {
		unixDate = time.Now().Unix()
	}

	ratesUntil, err := database.QueryAllExchangeRatesUntil(unixDate, client)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ratesUntilDto := models.ExchangeRatesForDateToDto(ratesUntil)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ratesUntilDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getRatesForCurrencyUntil(w http.ResponseWriter, r *http.Request, appCtx context.Context) {
	client, ok := appCtx.Value(database.MONGO_DB_CLIENT).(*mongo.Client)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Database Client not found")
		return
	}

	urlVars := mux.Vars(r)
	currency, ok := urlVars["currency"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("route-information missing: currency")
		return
	}

	unixDate, isFuture, err := parseDateAndHandleError(w, r, urlVars)
	if err != nil {
		return
	}

	if isFuture {
		unixDate = time.Now().Unix()
	}

	ratesUntil, err := database.QueryAllExchangeRatesUntil(unixDate, client)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ratesForCurrencyUntil := make([]models.ExchangeRateDto, len(ratesUntil))
	for _, entry := range ratesUntil {
		rate, ok := entry.ExchangeRates[currency]
		if ok {
			if parsedRate, err := strconv.ParseFloat(rate, 64); err == nil {
				date := time.Unix(entry.Date, 0).Format("2006-01-02")
				ratesForCurrencyUntil = append(ratesForCurrencyUntil, models.ExchangeRateDto{Currency: currency, Rate: parsedRate, Date: date})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ratesForCurrencyUntil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getRateForCurrency(w http.ResponseWriter, r *http.Request, appCtx context.Context) {
	client, ok := appCtx.Value(database.MONGO_DB_CLIENT).(*mongo.Client)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Database Client not found")
		return
	}

	urlVars := mux.Vars(r)
	currency, ok := urlVars["currency"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("route-information missing: currency")
		return
	}

	unixDate, isFuture, err := parseDateAndHandleError(w, r, urlVars)
	if err != nil {
		return
	}

	if isFuture {
		w.WriteHeader(http.StatusBadRequest)
		writeError(models.APIError{Message: "date is in the future", Code: models.FUTURE_DATE_ERROR}, w)
		return
	}

	rate, err := database.FindRateForCurrency(currency, unixDate, client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			handleDateNotFoundError(w, appCtx, currency, unixDate)
			return
		}
		if err == database.CurrencyError {
			w.WriteHeader(http.StatusNotFound)
			writeError(models.APIError{Message: err.Error(), Code: models.CURRENCY_NOT_FOUND_ERROR}, w)
		} // TODO: research MONGO-DB Errors

		w.WriteHeader(http.StatusNotFound)
		return
	}

	parsedRate, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeError(models.APIError{Message: "could not convert rate", Code: models.CONVERSION_ERROR}, w)
		return
	}

	date := time.Unix(unixDate, 0).Format("2006-01-02")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(models.ExchangeRateDto{Currency: currency, Rate: parsedRate, Date: date})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
