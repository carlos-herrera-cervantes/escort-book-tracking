package controllers

import (
    "context"
    "encoding/json"
    "net/http"
    "sync"

    "escort-book-tracking/config"
    "escort-book-tracking/models"
    "escort-book-tracking/repositories"
    "escort-book-tracking/services"
    "escort-book-tracking/types"

    "github.com/labstack/echo/v4"
    "github.com/labstack/gommon/log"
)

type CustomerTrackingController struct {
	Repository   repositories.ICustomerTrackingRepository
	KafkaService services.IKafkaService
}

func (h *CustomerTrackingController) GetCustomerLocation(c echo.Context) error {
	tracking, err := h.Repository.GetCustomerTracking(
	    c.Request().Context(),
	    c.Request().Header.Get("user-id"),
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, tracking)
}

func (h *CustomerTrackingController) SetCustomerLocation(c echo.Context) (err error) {
	var customerTracking models.CustomerTracking

	if err = c.Bind(&customerTracking.Location); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	userId := c.Request().Header.Get("user-id")
    customerTracking.CustomerId = userId
	ctx := c.Request().Context()

	if err = h.Repository.AlterCustomerTracking(ctx, &customerTracking); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go h.emitAcknowledge(ctx, userId, &wg)

	wg.Wait()

	return c.JSON(http.StatusCreated, customerTracking)
}

func (h *CustomerTrackingController) emitAcknowledge(ctx context.Context, userId string, wg *sync.WaitGroup) {
	defer wg.Done()
	tracking, err := h.Repository.GetCustomerTracking(ctx, userId)

	if err != nil {
		log.Errorf("customer_controller->emitAcknowledge->GetCustomerTracking:%s", err.Error())
		return
	}

	if tracking.Acknowledged {
		log.Infof("customer_controller->emitAcknowledge->The customer %s already acknowledged", userId)
		return
	}

	if err := h.Repository.Acknowledge(ctx, userId); err != nil {
		log.Errorf("customer_controller->emitAcknowledge->Acknowledge:%s", err.Error())
		return
	}

	countUserEvent := types.CountUserEvent{
		Accumulator: 1,
		Operation:   config.InitOperationConfig().NewUser,
		UserId:      userId,
		UserType:    "Customer",
	}
	bytes, _ := json.Marshal(countUserEvent)

	if err := h.KafkaService.SendMessage(config.InitKafkaConfig().Topics.OperationTopic, bytes); err != nil {
		log.Errorf("customer_controller->emitAcknowledge->SendMessage:%s", err.Error())
	}
}
