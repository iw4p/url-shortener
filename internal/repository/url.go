package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type URLRepository struct {
	collection *mongo.Collection
}

func NewURLRepository(collection *mongo.Collection) *URLRepository {
	return &URLRepository{collection: collection}
}

func (r *URLRepository) GetLastShortValue(ctx context.Context) (int64, error) {
	opts := options.FindOne().SetSort(bson.D{{"short_id", -1}})
	var result struct {
		ShortId int64 `bson:"short_id"`
	}
	err := r.collection.FindOne(ctx, bson.D{}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}
	return result.ShortId, nil
}

func (r *URLRepository) InsertDocument(ctx context.Context, data bson.D) (interface{}, error) {
	res, err := r.collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func (r *URLRepository) GetDocument(ctx context.Context, filter bson.D) (bson.M, error) {
	var result bson.M
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
