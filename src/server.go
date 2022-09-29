package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

type handlerWithAppContext = func(http.ResponseWriter, *http.Request, context.Context)

func startWebSever(appCtx context.Context) {
	r := mux.NewRouter()
	r.HandleFunc("/rate/{currency}/{date}", useAppContext(getRateForCurrency, appCtx))
	r.HandleFunc("/supported", useAppContext(getSupportedCurrencies, appCtx))
	r.HandleFunc("/supported/{date}", useAppContext(getSupportedCurrencies, appCtx))
	http.Handle("/", r)
	fmt.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":5555", nil))
}

func useAppContext(next handlerWithAppContext, appCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, appCtx)
	}
}

func getSupportedCurrencies(w http.ResponseWriter, r *http.Request, appCtx context.Context) {
	client, ok := appCtx.Value(MONGO_DB_CLIENT).(*mongo.Client)
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
			writeError(APIError{"date format not supported", FORMAT_NOT_SUPPORTED_ERROR}, w)
			return
		}
		today := time.Now()
		isFuture := today.Before(parsedDate)
		if isFuture {
			w.WriteHeader(http.StatusBadRequest)
			writeError(APIError{"date is in the future", FUTURE_DATE_ERROR}, w)
			return
		}
		unixDate = parsedDate.Unix()
	}

	supportedCurrencies, err := findAllSupportedCurrencies(unixDate, client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			writeError(APIError{"could not find any information", NO_ENTRY_FOUND_ERROR}, w)
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

func getRateForCurrency(w http.ResponseWriter, r *http.Request, appCtx context.Context) {
	client, ok := appCtx.Value(MONGO_DB_CLIENT).(*mongo.Client)
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
	date, ok := urlVars["date"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("route-information missing: date")
		return
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeError(APIError{"date format not supported", FORMAT_NOT_SUPPORTED_ERROR}, w)
		// TODO: Use error structs with more detail
		return
	}
	unixDate := parsedDate.Unix()

	today := time.Now()
	isFuture := today.Before(parsedDate)

	if isFuture {
		w.WriteHeader(http.StatusBadRequest)
		writeError(APIError{"date is in the future", FUTURE_DATE_ERROR}, w)
		return
	}

	rate, err := findRateForCurrency(currency, parsedDate.Unix(), client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			handleDateNotFoundError(w, appCtx, currency, unixDate)
			return
		}
		if err == CurrencyError {
			w.WriteHeader(http.StatusNotFound)
			writeError(APIError{err.Error(), CURRENCY_NOT_FOUND_ERROR}, w)
		} // TODO: research MONGO-DB Errors

		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ExchangeRate{currency, rate, unixDate})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func handleDateNotFoundError(w http.ResponseWriter, appCtx context.Context, currency string, unixDate int64) {
	err := refreshTodaysRates(appCtx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client, ok := appCtx.Value(MONGO_DB_CLIENT).(*mongo.Client)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("database client not found") // TODO: HANDLE HANDLE HANDLE
		return
	}
	rate, err := findRateForCurrency(currency, unixDate, client)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		writeError(APIError{"could not get rate for this currency and unixDate", NO_ENTRY_FOUND_ERROR}, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ExchangeRate{currency, rate, unixDate})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func writeError(error APIError, w http.ResponseWriter) {
	err := json.NewEncoder(w).Encode(error)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
