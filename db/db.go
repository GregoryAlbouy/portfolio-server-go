package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"gregoryalbouy-server-go/utl"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// Client reference
	Client *mongo.Client
	// DB reference
	DB *mongo.Database
)

// Connect sets Client from the given URI and DB from given name
// if it exists in the Client
func Connect(uri, name string) {
	defer func(t time.Time) {
		fmt.Println("DB connection time: ", time.Since(t))
	}(time.Now())

	Client, DB = initDB(uri, name)
}

// Collection returns queried collection from current database
func Collection(name string) *mongo.Collection {
	return DB.Collection(name)
}

func initDB(uri, name string) (client *mongo.Client, db *mongo.Database) {
	const timeout = 10 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	utl.Check(err)

	err = client.Ping(ctx, nil)
	utl.Check(err)

	db = client.Database(name)

	return
}

func ping(ctx context.Context, client *mongo.Client) {
	fmt.Println("Sending Ping...")
	err := client.Ping(ctx, nil)
	utl.Check(err)
	// fmt.Println("Ping ok!")
}

func printDatabases(ctx context.Context, client *mongo.Client) {
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
}
