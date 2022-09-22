package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const MONGO_DB_CLIENT = "mongodb-client"
const DB_ECB_RATES = "ecb-rates-db"
const COL_EX_RATES = "exchange-rates"

func getDatabaseClient(appCtx context.Context) *mongo.Client {
	return appCtx.Value(MONGO_DB_CLIENT).(*mongo.Client)
}

func useDatabase(child func(appCtx context.Context), appCtx context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	appCtx = context.WithValue(appCtx, MONGO_DB_CLIENT, client)

	child(appCtx)

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func saveExchangeRatesToDB(exchangeRates *ExchangeRatesForDate, client *mongo.Client) {

	println("Saving rates to db...")
	collection := client.Database(DB_ECB_RATES).Collection(COL_EX_RATES)

	_, err := collection.InsertOne(context.TODO(), exchangeRates)
	if err != nil {
		fmt.Errorf("Could not save rates due to error: %v\n", err)
	}
}
