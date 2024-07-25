package service

import (
	"context"
	"fmt"
	"time"

	"github.com/iw4p/url-shortener/base62"
	"github.com/iw4p/url-shortener/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
)

type URLService struct {
	repo URLRepository
}

type URLRequest struct {
	Short    string `json:"short"`
	Original string `json:"original"`
}

type URLResponse struct {
	Short    string `json:"short"`
	Original string `json:"original"`
}

type URLRepository interface {
	GetLastShortValue(ctx context.Context) (int64, error)
	InsertDocument(ctx context.Context, data bson.D) (*repository.Document, error)
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
	result, err := s.repo.InsertDocument(ctx, data)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	urlResponse := URLResponse{
		Short:    result.Short,
		Original: result.Original,
	}

	return urlResponse, nil
}

func (s *URLService) GetOriginal(ctx context.Context, short string) (*repository.Document, error) {

	b := base62.Base62{}.DecodeBase62(short)
	filter := bson.D{{"short_id", b}}
	result, err := s.repo.GetDocument(ctx, filter)

	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	short, ok := result["short"].(string)
	if !ok {
		return nil, fmt.Errorf("short id not found or not an int32")
	}

	original, ok := result["original"].(string)

	if !ok || original == "" {
		return nil, fmt.Errorf("url is not available")
	}

	urlResponse := repository.Document{
		Short:    short,
		Original: original,
	}

	return &urlResponse, nil
}
