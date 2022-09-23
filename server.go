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

func startWebSever(appCtx context.Context) {
	r := mux.NewRouter()
	r.Handle("/rate/{currency}/{date}", useAppContext(http.HandlerFunc(getRateForCurrency), appCtx))
	http.Handle("/", r)
	fmt.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":5555", nil))
}

func useAppContext(next http.Handler, appCtx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(appCtx))
	})
}

func getRateForCurrency(w http.ResponseWriter, r *http.Request) {
	client, ok := r.Context().Value(MONGO_DB_CLIENT).(*mongo.Client)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Database Client not found")
		return
	}

	urlVars := mux.Vars(r)
	currency, ok := urlVars["currency"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	date, ok := urlVars["date"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rate, err := findRateForCurrency(currency, date, client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			handleDateNotFoundError(w, r)
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
	err = json.NewEncoder(w).Encode(ExchangeRate{currency, rate, date})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleDateNotFoundError(w http.ResponseWriter, r *http.Request) {
	urlVars := mux.Vars(r)
	currency := urlVars["currency"]
	date := urlVars["date"]

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("date format not supported"))
		// TODO: Use error structs with more detail
		return
	}

	today := time.Now()
	isFuture := today.Before(parsedDate)

	if isFuture {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("date is in the future"))
		return
	}
	err = refreshTodaysRates(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client, ok := r.Context().Value(MONGO_DB_CLIENT).(*mongo.Client)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("database client not found") // TODO: HANDLE HANDLE HANDLE
		return
	}
	rate, err := findRateForCurrency(currency, date, client)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("could not get rate for this currency and date"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ExchangeRate{currency, rate, date})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
