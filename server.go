package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func getRateForCurrency(date string, currency string, appCtx context.Context) string {
	client := getDatabaseClient(appCtx)
	collection := client.Database(DB_ECB_RATES).Collection(COL_EX_RATES)

	var rates ExchangeRatesForDate
	err := collection.FindOne(context.TODO(), bson.D{{"date", date}}).Decode(&rates)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// TODO: ERROR 404 NOT FOUND
			return ""
		}
		log.Printf("Collection error: %v\n", err)
		return ""
	}
	rate, ok := rates.ExchangeRates[currency]
	if ok {
		return rate
	}
	return "" // TODO: NOT FOUND
}
