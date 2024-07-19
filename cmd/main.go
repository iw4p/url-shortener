package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/iw4p/url-shortener/utils"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Response struct {
	OK      bool
	Message string
}

type shortURLRequest struct {
	Url string `json:"url"`
}

func healthCheck(c echo.Context) error {
	res := Response{
		OK:      true,
		Message: "heartbeat is ok",
	}
	return c.JSON(http.StatusOK, res)
}

func shortURL(c echo.Context) error {
	r := new(shortURLRequest)
	if err := c.Bind(r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

func main() {

	environmentVariables := utils.ReturnEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(environmentVariables.Mongo_url))
	if err != nil {
		panic("can not connect to the DB")
	}

	collection := client.Database(environmentVariables.Db_name).Collection(environmentVariables.Collection)

	fmt.Println(collection.Indexes())

	// Query the collection using the indexed field
	doc, err := getDataByID(ctx, collection, "1231231ww2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found document: %+v\n", doc)

	// Insert a document
	docID, err := insertDocument(ctx, collection)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Inserted document with ID: %v\n", docID)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("/api/v1")

	v1.GET("/health", healthCheck)
	v1.POST("/short", shortURL)

	e.Logger.Fatal(e.Start(":8080"))
}

func insertDocument(ctx context.Context, collection *mongo.Collection) (interface{}, error) {
	res, err := collection.InsertOne(ctx, bson.D{
		{"short", "1231231ww22"},
		{"original", "https://google.com"},
	})
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func getDataByID(ctx context.Context, collection *mongo.Collection, short string) (bson.M, error) {
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
