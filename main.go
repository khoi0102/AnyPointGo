package main

import (
	"AnyPointGo/racer"
	"fmt"
	"github.com/labstack/echo/v4"
	"os"
)

func main() {
	e := echo.New()
	e.POST("/races", racer.Race)
	e.GET("/", racer.Hello)
	e.POST("/races/:id/laps", racer.RaceLap)
	e.POST("/temperatures", racer.Temperature)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
