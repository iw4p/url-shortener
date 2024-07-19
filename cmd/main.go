package main

import (
	"net/http"

	"github.com/labstack/echo"
)

type Response struct {
	OK      bool
	Message string
}

func healthCheck(c echo.Context) error {
	res := Response{
		OK:      true,
		Message: "heartbeat is ok",
	}
	return c.JSON(http.StatusOK, res)
}

func main() {

	e := echo.New()
	e.GET("/health", healthCheck)

	e.Logger.Fatal(e.Start(":8080"))
}
