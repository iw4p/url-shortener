package service

import (
	"context"
	"time"

	"github.com/iw4p/url-shortener/base62"
	"go.mongodb.org/mongo-driver/bson"
)

type URLService struct {
	repo URLRepository
}

type URLRequest struct {
	Short    string `json:"short"`
	Original string `json:"original"`
}

type URLRepository interface {
	GetLastShortValue(ctx context.Context) (int64, error)
	InsertDocument(ctx context.Context, data bson.D) (interface{}, error)
	GetDocument(ctx context.Context, filter bson.D) (bson.M, error)
}

func NewURLService(repo URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) GetShorten(ctx context.Context, original string) (interface{}, error) {
	lastShort, err := s.repo.GetLastShortValue(ctx)

	if err != nil {
		return nil, err
	}

	var nextSeq int64
	if lastShort == 0 {
		nextSeq = 1
	} else {
		nextSeq = lastShort + 1
	}
	shortValue := base62.Base62{}.EncodeBase62(nextSeq)

	data := bson.D{{"short_id", nextSeq}, {"short", shortValue}, {"original", original}, {"createdAt", time.Now()}}
	return s.repo.InsertDocument(ctx, data)
}

func (s *URLService) GetOriginal(ctx context.Context, short string) (bson.M, error) {

	b := base62.Base62{}.DecodeBase62(short)
	filter := bson.D{{"short_id", b}}
	return s.repo.GetDocument(ctx, filter)
}
