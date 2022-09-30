package database

import (
	"context"
	"damattl.de/api/currency/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

const (
	DEFAULT_MONGO_URL = "mongodb://localhost:27017"
	MONGO_DB_CLIENT   = "mongodb-client"
	DB_ECB_RATES      = "linum_exchange_rates_db"
	COL_EX_RATES      = "exchange-rates"
)

func GetClient(appCtx context.Context) *mongo.Client {
	return appCtx.Value(MONGO_DB_CLIENT).(*mongo.Client)
}

func getDatabaseURL() string {
	url := os.Getenv("MONGO_URL")
	fmt.Println(url)
	if url == "" {
		fmt.Println("USING DEFAULT URL")
		return DEFAULT_MONGO_URL
	}
	return url
}

func UseDatabase(child func(appCtx context.Context), appCtx context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(getDatabaseURL()))
	if err != nil {
		panic("can't connect to database")
	}

	appCtx = context.WithValue(appCtx, MONGO_DB_CLIENT, client)

	child(appCtx)

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func SaveExchangeRates(exchangeRates *models.ExchangeRatesForDate, client *mongo.Client) {

	println("Saving rates to db...")
	collection := client.Database(DB_ECB_RATES).Collection(COL_EX_RATES)

	_, err := collection.InsertOne(context.TODO(), exchangeRates)
	if err != nil {
		fmt.Errorf("Could not save rates due to error: %v\n", err)
	}
}

func FindRateForCurrency(currency string, date int64, client *mongo.Client) (string, error) {
	collection := client.Database(DB_ECB_RATES).Collection(COL_EX_RATES)

	var rates models.ExchangeRatesForDate
	err := collection.FindOne(context.TODO(), bson.D{{"date", date}}).Decode(&rates)
	if err != nil {
		return "", err
	}
	rate, ok := rates.ExchangeRates[currency]
	if !ok {
		return "", CurrencyError
	}
	return rate, nil
}

func FindAllSupportedCurrencies(date int64, client *mongo.Client) ([]string, error) {
	collection := client.Database(DB_ECB_RATES).Collection(COL_EX_RATES)
	var rates models.ExchangeRatesForDate
	if date == -1 {
		opts := options.FindOne().SetSort(bson.M{"$natural": -1})
		err := collection.FindOne(context.TODO(), bson.M{}, opts).Decode(&rates)
		if err != nil {
			return nil, err
		}
	} else {
		err := collection.FindOne(context.TODO(), bson.D{{"date", date}}).Decode(&rates)
		if err != nil {
			return nil, err
		}
	}
	supportedCurrencies := make([]string, 0, len(rates.ExchangeRates))
	for key := range rates.ExchangeRates {
		supportedCurrencies = append(supportedCurrencies, key)
	}
	return supportedCurrencies, nil
}

func QueryAllExchangeRatesUntil(date int64, client *mongo.Client) ([]models.ExchangeRatesForDate, error) {
	collection := client.Database(DB_ECB_RATES).Collection(COL_EX_RATES)
	var ratesUntil []models.ExchangeRatesForDate
	filter := bson.M{
		"date": bson.M{
			"$lte": date,
		},
	}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &ratesUntil); err != nil {
		return nil, err
	}
	return ratesUntil, nil
}
