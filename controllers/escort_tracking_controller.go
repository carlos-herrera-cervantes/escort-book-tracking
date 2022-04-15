package controllers

import (
	"escort-book-tracking/models"
	"escort-book-tracking/repositories"
	"escort-book-tracking/types"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type EscortTrackingController struct {
	Repository              *repositories.EscortTrackingRepository
	EscortProfileRepository *repositories.EscortProfileRepository
}

func (h *EscortTrackingController) GetLocationsByTerritory(c echo.Context) (err error) {
	var pager types.Pager

	c.Bind(&pager)

	if err = pager.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	territory := c.QueryParam("territory")
	trackings, _ := h.Repository.GetEscortLocationByTerritory(c.Request().Context(), territory, pager.Offset, pager.Limit)

	for index, value := range trackings {
		profile, _ := h.EscortProfileRepository.GetEscortProfile(c.Request().Context(), value.EscortId)
		trackings[index].FirstName = profile.FirstName
		trackings[index].LastName = profile.LastName
		trackings[index].Avatar = fmt.Sprintf("%s/%s/%s", os.Getenv("ENDPOINT"), os.Getenv("S3"), profile.Avatar)
	}

	number, _ := h.Repository.CountEscortLocationByTerritory(c.Request().Context())
	pagerResult := types.PagerResult{}

	return c.JSON(http.StatusOK, pagerResult.GetPagerResult(&pager, number, trackings))
}

func (h *EscortTrackingController) GetEscortLocation(c echo.Context) error {
	tracking, err := h.Repository.GetEscortTracking(c.Request().Context(), c.Request().Header.Get("user-id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, tracking)
}

func (h *EscortTrackingController) SetEscortLocation(c echo.Context) (err error) {
	var escortTracking models.EscortTracking
	escortTracking.EscortId = c.Request().Header.Get("user-id")

	if err = c.Bind(&escortTracking.Location); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err = h.Repository.UpsertEscortTracking(c.Request().Context(), &escortTracking); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, escortTracking)
}
