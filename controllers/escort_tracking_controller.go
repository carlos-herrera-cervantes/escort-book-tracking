package controllers

import (
	"escort-book-tracking/models"
	"escort-book-tracking/repositories"
	"escort-book-tracking/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

type EscortTrackingController struct {
	Repository *repositories.EscortTrackingRepository
}

func (h *EscortTrackingController) GetLocationsByTerritory(c echo.Context) (err error) {
	var pager types.Pager

	c.Bind(&pager)

	if err = pager.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	territory := c.QueryParam("territory")
	trackings, _ := h.Repository.GetByTerritory(c.Request().Context(), territory, pager.Offset, pager.Limit)
	number, _ := h.Repository.Count(c.Request().Context())
	pagerResult := types.PagerResult{}

	return c.JSON(http.StatusOK, pagerResult.GetPagerResult(&pager, number, trackings))
}

func (h *EscortTrackingController) GetEscortLocation(c echo.Context) error {
	escortId := ""
	tracking, err := h.Repository.GetOne(c.Request().Context(), escortId)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, tracking)
}

func (h *EscortTrackingController) SetEscortLocation(c echo.Context) (err error) {
	var escortTracking models.EscortTracking
	escortId := ""
	escortTracking.EscortId = escortId

	if err = c.Bind(&escortTracking.Location); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err = h.Repository.UpsertOne(c.Request().Context(), &escortTracking); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, escortTracking)
}
