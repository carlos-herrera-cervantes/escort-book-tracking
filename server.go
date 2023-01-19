package main

import (
    "fmt"

    "escort-book-tracking/config"
    "escort-book-tracking/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	v1 := e.Group("/api/v1")

	routes.BootstrapCustomerTrackingRoutes(v1)
	routes.BootstrapEscortTrackingRoutes(v1)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", config.InitApp().Port)))
}
