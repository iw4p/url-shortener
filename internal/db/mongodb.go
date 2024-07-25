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

	mod := mongo.IndexModel{
		Keys: bson.M{
			"short_id": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	ind, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		return nil, nil, fmt.Errorf("Indexes().CreateOne() ERROR:%w", err)
	} else {
		fmt.Println("CreateOne() index:", ind)
	}

	return ctx, collection, nil
}

func (m *MongoDB) Disconnect(ctx context.Context) error {
	if m.client == nil {
		return nil
	}
	return m.client.Disconnect(ctx)
}
