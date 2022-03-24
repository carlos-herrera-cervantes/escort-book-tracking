package routes

import (
	"escort-book-tracking/controllers"
	"escort-book-tracking/db"
	"escort-book-tracking/repositories"

	"github.com/labstack/echo/v4"
)

func BoostrapEscortTrackingRoutes(v *echo.Group) {
	router := &controllers.EscortTrackingController{
		Repository: &repositories.EscortTrackingRepository{
			Data: db.New(),
		},
	}

	v.GET("/tracking/escort", router.GetEscortLocation)
	v.GET("/tracking/escorts", router.GetLocationsByTerritory)
	v.POST("/tracking/escort", router.SetEscortLocation)
}
