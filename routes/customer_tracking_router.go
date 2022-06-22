package routes

import (
	"escort-book-tracking/controllers"
	"escort-book-tracking/db"
	"escort-book-tracking/repositories"
	"escort-book-tracking/services"

	"github.com/labstack/echo/v4"
)

func BoostrapCustomerTrackingRoutes(v *echo.Group) {
	router := &controllers.CustomerTrackingController{
		Repository: &repositories.CustomerTrackingRepository{
			Data: db.InitDB("default"),
		},
		KafkaService: &services.KafkaService{
			Producer: db.NewProducer(),
		},
	}

	v.GET("/tracking/customer", router.GetCustomerLocation)
	v.POST("/tracking/customer", router.SetCustomerLocation)
}
