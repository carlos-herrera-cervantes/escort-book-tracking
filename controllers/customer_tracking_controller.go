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
	tracking, err := h.Repository.GetCustomerTracking(c.Request().Context(), c.Request().Header.Get("user-id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, tracking)
}

func (h *CustomerTrackingController) SetCustomerLocation(c echo.Context) (err error) {
	var customerTracking models.CustomerTracking
	customerTracking.CustomerId = c.Request().Header.Get("user-id")

	if err = c.Bind(&customerTracking.Location); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err = h.Repository.UpsertCustomerTracking(c.Request().Context(), &customerTracking); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, customerTracking)
}
