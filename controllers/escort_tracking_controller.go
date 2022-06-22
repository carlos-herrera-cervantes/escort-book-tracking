package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"escort-book-tracking/config"
	"escort-book-tracking/models"
	"escort-book-tracking/repositories"
	"escort-book-tracking/services"
	"escort-book-tracking/types"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type EscortTrackingController struct {
	Repository              repositories.IEscortTrackingRepository
	EscortProfileRepository repositories.IEscortProfileRepository
	KafkaService            services.IKafkaService
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

	var wg sync.WaitGroup
	wg.Add(1)

	go h.emitAcknowledge(c.Request().Context(), c.Request().Header.Get("user-id"), &wg)

	wg.Wait()

	return c.JSON(http.StatusCreated, escortTracking)
}

func (h *EscortTrackingController) emitAcknowledge(ctx context.Context, userId string, wg *sync.WaitGroup) {
	defer wg.Done()
	tracking, err := h.Repository.GetEscortTracking(ctx, userId)

	if err != nil {
		log.Errorf("escort_controller->emitAcknowledge->GetEscortTracking:%s", err.Error())
		return
	}

	if tracking.Acknowledged {
		log.Infof("escort_controller->emitAcknowledge->The escort %s already acknowledged", userId)
		return
	}

	if err := h.Repository.Acknowledge(ctx, userId); err != nil {
		log.Errorf("escort_controller->emitAcknowledge->Acknowledge:%s", err.Error())
		return
	}

	countUserEvent := types.CountUserEvent{
		Accumulator: 1,
		Operation:   config.NewUser,
		UserId:      userId,
		UserType:    "Escort",
	}
	bytes, _ := json.Marshal(countUserEvent)

	if err := h.KafkaService.SendMessage(ctx, config.OperationTopic, bytes); err != nil {
		log.Errorf("escort_controller->emitAcknowledge->SendMessage:%s", err.Error())
	}
}
