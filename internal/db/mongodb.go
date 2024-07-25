package db

import (
	"context"
	"fmt"
	"time"

	"github.com/iw4p/url-shortener/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
}

func NewMongoDB() *MongoDB {
	return &MongoDB{}
}

func (m *MongoDB) Init(collectionName string) (context.Context, *mongo.Collection, error) {
	env := config.ReturnEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(env.MongoURL))
	if err != nil {
		return nil, nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, nil, err
	}

	m.client = client
	collection := client.Database(env.DbName).Collection(collectionName)

	// Create unique index on short_id
	uniqueIndexModel := mongo.IndexModel{
		Keys: bson.M{
			"short_id": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	// Create TTL index on createdAt
	ttlIndexModel := mongo.IndexModel{
		Keys: bson.M{
			"createdAt": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(157_800_000), // Remove each URL records which passed 5 Years (157,800,000 seconds)
	}

	_, err = collection.Indexes().CreateOne(ctx, uniqueIndexModel)
	if err != nil {
		fmt.Println("Indexes().CreateOne() ERROR:", err)
		return nil, nil, err
	}

	_, err = collection.Indexes().CreateOne(ctx, ttlIndexModel)
	if err != nil {
		fmt.Println("Indexes().CreateOne() ERROR:", err)
		return nil, nil, err
	}

	return ctx, collection, nil
}

func (m *MongoDB) Disconnect(ctx context.Context) error {
	if m.client == nil {
		return nil
	}
	return m.client.Disconnect(ctx)
}
