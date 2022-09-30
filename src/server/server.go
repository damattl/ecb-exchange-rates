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
	r := mux.NewRouter()
	r.HandleFunc("/rate/{date}/{currency}", useAppContext(getRateForCurrency, appCtx))
	r.HandleFunc("/rates-until/{date}", useAppContext(getRatesUntil, appCtx))
	r.HandleFunc("/rates-until/{date}/{currency}", useAppContext(getRatesForCurrencyUntil, appCtx))
	r.HandleFunc("/supported", useAppContext(getSupportedCurrencies, appCtx))
	r.HandleFunc("/supported/{date}", useAppContext(getSupportedCurrencies, appCtx))
	http.Handle("/", r)
	fmt.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":5555", nil))
}
