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
	http.Handle("/", r)
	fmt.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":5555", nil))
}

func useAppContext(next handlerWithAppContext, appCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, appCtx)
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
		w.Write([]byte("date format not supported"))
		// TODO: Use error structs with more detail
		return
	}
	unixDate := parsedDate.Unix()

	today := time.Now()
	isFuture := today.Before(parsedDate)

	if isFuture {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("date is in the future"))
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
			w.Write([]byte(err.Error()))
		}

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
		w.Write([]byte("could not get rate for this currency and unixDate"))
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
