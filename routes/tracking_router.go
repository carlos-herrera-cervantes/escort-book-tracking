package routes

import "github.com/labstack/echo/v4"

func BoostrapTrackingRoutes(v *echo.Group) {
	v.GET("/tracking/customer", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "OK",
		})
	})
}
