package main

import (
	"fmt"
	"net/http"

	"github.com/iw4p/url-shortener/base62"
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
	b := base62.Base62{}

	fmt.Println(b.DecodeBase62("adwadwa"))
	fmt.Println(b.EncodeBase62(1232131233))

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("/api/v1")

	v1.GET("/health", healthCheck)
	v1.POST("/short", shortURL)

	e.Logger.Fatal(e.Start(":8080"))
}
