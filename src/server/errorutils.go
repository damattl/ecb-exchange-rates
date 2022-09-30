package server

import (
	"context"
	"damattl.de/api/currency/database"
	"damattl.de/api/currency/models"
	"damattl.de/api/currency/tasks"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"time"
)

func handleDateNotFoundError(w http.ResponseWriter, appCtx context.Context, currency string, unixDate int64) {
	err := tasks.RefreshTodaysRates(appCtx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client, ok := appCtx.Value(database.MONGO_DB_CLIENT).(*mongo.Client)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("database client not found") // TODO: HANDLE HANDLE HANDLE
		return
	}
	rate, err := database.FindRateForCurrency(currency, unixDate, client)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		writeError(models.APIError{Message: "could not get rate for this currency and unixDate", Code: models.NO_ENTRY_FOUND_ERROR}, w)
		return
	}

	parsedRate, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeError(models.APIError{Message: "could not convert rate", Code: models.CONVERSION_ERROR}, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(models.ExchangeRate{Currency: currency, Rate: parsedRate, Date: unixDate})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func writeError(error models.APIError, w http.ResponseWriter) {
	err := json.NewEncoder(w).Encode(error)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func parseDateAndHandleError(w http.ResponseWriter, r *http.Request, urlVars map[string]string) (int64, bool, error) {
	date, ok := urlVars["date"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return 0, false, errors.New("route-information missing: date")
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeError(models.APIError{Message: "date format not supported", Code: models.FORMAT_NOT_SUPPORTED_ERROR}, w)
		// TODO: Use error structs with more detail
		return 0, false, err
	}
	unixDate := parsedDate.Unix()

	today := time.Now()
	isFuture := today.Before(parsedDate)

	return unixDate, isFuture, nil
}
