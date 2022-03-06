package main

import (
	"escort-book-tracking/routes"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	v1 := e.Group("/api/v1")

	routes.BoostrapTrackingRoutes(v1)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", os.Getenv("PORT"))))
}
