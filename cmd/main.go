package main

import (
	"log"

	"github.com/iw4p/url-shortener/config"
	"github.com/iw4p/url-shortener/handler"
	"github.com/iw4p/url-shortener/internal/db"
	"github.com/iw4p/url-shortener/internal/repository"
	"github.com/iw4p/url-shortener/internal/service"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	environmentVariables := config.ReturnEnv()

	mongoDB := db.NewMongoDB()

	ctx, collection, err := mongoDB.Init(environmentVariables.Collection)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer mongoDB.Disconnect(ctx)

	urlRepository := repository.NewURLRepository(collection)
	urlService := service.NewURLService(urlRepository)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := handler.NewHandler(urlService)

	e.GET("/:redirect", h.Redirect)

	v1 := e.Group("/api/v1")
	v1.GET("/health", h.HealthCheck)
	v1.POST("/short", h.ShortURL)
	v1.POST("/original", h.GetOriginalURL)

	e.Logger.Fatal(e.Start(":8080"))
}
