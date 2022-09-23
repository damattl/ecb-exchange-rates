package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
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

	rate, ok := findRateForCurrency(currency, date, client)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(ExchangeRate{currency, rate, date})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
