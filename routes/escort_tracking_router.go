package routes

import (
	"escort-book-tracking/controllers"
	"escort-book-tracking/db"
	"escort-book-tracking/repositories"
	"escort-book-tracking/services"

	"github.com/labstack/echo/v4"
)

func BootstrapEscortTrackingRoutes(v *echo.Group) {
	router := &controllers.EscortTrackingController{
		Repository: &repositories.EscortTrackingRepository{
			Data: db.NewPostgresClient(),
		},
		EscortProfileRepository: &repositories.EscortProfileRepository{
			Data: db.NewPostgresClient(),
		},
		KafkaService: &services.KafkaService{
			Producer: db.NewProducer(),
		},
	}

	v.GET("/tracking/escort", router.GetEscortLocation)
	v.GET("/tracking/escorts", router.GetLocationsByTerritory)
	v.POST("/tracking/escort", router.SetEscortLocation)
}
