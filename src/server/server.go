package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type handlerWithAppContext = func(http.ResponseWriter, *http.Request, context.Context)

func StartWebSever(appCtx context.Context) {
	dateRegex := "[1-9][0-9]{3}-(?:0[1-9]|1[0-2])-(?:0[1-9]|[12][0-9]|3[0-1])"
	currencyRegex := "[a-zA-Z]{3}"

	r := mux.NewRouter()
	r.HandleFunc(fmt.Sprintf("/rate/{date:%s}/{currency:%s}", dateRegex, currencyRegex), useAppContext(getRateForCurrency, appCtx))
	r.HandleFunc(fmt.Sprintf("/rates/{date:%s}", dateRegex), useAppContext(getRatesForDate, appCtx))
	r.HandleFunc(fmt.Sprintf("/rates-until/{date:%s}", dateRegex), useAppContext(getRatesUntil, appCtx))
	r.HandleFunc(fmt.Sprintf("/rates-until/{date:%s}/{currency:%s}", dateRegex, currencyRegex), useAppContext(getRatesForCurrencyUntil, appCtx))
	r.HandleFunc(fmt.Sprintf("/rates-between/{earliestDate:%s}/{latestDate:%s}", dateRegex, dateRegex), useAppContext(getRatesForTimeSpan, appCtx))
	r.HandleFunc("/supported", useAppContext(getSupportedCurrencies, appCtx))
	r.HandleFunc(fmt.Sprintf("/supported/{date:%s}", dateRegex), useAppContext(getSupportedCurrencies, appCtx))
	http.Handle("/", r)
	fmt.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":5555", nil))
}
