package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertDocument(ctx context.Context, collection *mongo.Collection) (interface{}, error) {
	res, err := collection.InsertOne(ctx, bson.D{
		{"short", "1231231ww22"},
		{"original", "https://google.com"},
	})
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func GetDataByID(ctx context.Context, collection *mongo.Collection, short string) (bson.M, error) {
	filter := bson.D{{"short", short}}

	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no document found with the given id")
		}
		return nil, err
	}

	return result, nil
}
