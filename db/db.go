package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func GetDBCollection(coll string) *mongo.Collection {
	return db.Collection(coll)
}

func Init() error {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal(
			"Please set MONGODB_URI in your .env",
		)
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	db = client.Database("Sangeet")
	return nil
}

func Close() error {
	if err := db.Client().Disconnect(context.TODO()); err != nil {
		return err
	}
	return nil
}
