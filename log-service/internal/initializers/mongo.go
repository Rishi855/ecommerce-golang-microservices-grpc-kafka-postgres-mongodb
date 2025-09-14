package initializers

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var LogCollection *mongo.Collection

func ConnectMongoDB(uri, dbName, collectionName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("[Log-Service] Failed to connect MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("[Log-Service] MongoDB ping failed: %v", err)
	}

	log.Println("[Log-Service] Connected to MongoDB")

	MongoClient = client
	LogCollection = client.Database(dbName).Collection(collectionName)
}
