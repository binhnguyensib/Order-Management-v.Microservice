package config

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func ConnectMongoDB() (*MongoDB, error) {
	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	db := client.Database(dbName)
	return &MongoDB{
		Client: client,
		DB:     db}, nil

}

func (d *MongoDB) Close() error {
	if err := d.Client.Disconnect(context.TODO()); err != nil {
		log.Printf("Failed to disconnect from MongoDB: %v", err)
		return err
	}
	return nil
}
