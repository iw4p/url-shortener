package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/iw4p/url-shortener/db"
	"github.com/iw4p/url-shortener/repository"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

	mongoDB := db.MongoDB{}

	ctx, collection, err := mongoDB.InitMongoDB("urls")
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}

	// Query the collection using the indexed field
	doc, err := repository.GetDataByID(ctx, collection, "1231231ww2")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found document: %+v\n", doc)

	// Insert a document
	docID, err := repository.InsertDocument(ctx, collection)
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
