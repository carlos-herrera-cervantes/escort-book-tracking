package controllers

import (
	"escort-book-tracking/models"
	"escort-book-tracking/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CustomerTrackingController struct {
	Repository *repositories.CustomerTrackingRepository
}

func (h *CustomerTrackingController) GetCustomerLocation(c echo.Context) error {
	customerId := ""
	tracking, err := h.Repository.GetOne(c.Request().Context(), customerId)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, tracking)
}

func (h *CustomerTrackingController) SetCustomerLocation(c echo.Context) (err error) {
	var customerTracking models.CustomerTracking
	customerId := ""
	customerTracking.CustomerId = customerId

	if err = c.Bind(&customerTracking.Location); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err = h.Repository.UpsertOne(c.Request().Context(), &customerTracking); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, customerTracking)
}
