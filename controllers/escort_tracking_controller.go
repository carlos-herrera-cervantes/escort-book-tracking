package controllers

import (
    "context"
    "encoding/json"
    "fmt"
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

type EscortTrackingController struct {
	Repository              repositories.IEscortTrackingRepository
	EscortProfileRepository repositories.IEscortProfileRepository
	KafkaService            services.IKafkaService
}

func (h *EscortTrackingController) GetLocationsByTerritory(c echo.Context) (err error) {
	pager := types.Pager{}

	if err := c.Bind(&pager); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

	if err = pager.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	territory := c.QueryParam("territory")
	ctx := c.Request().Context()
	tracking, err := h.Repository.GetEscortLocationByTerritory(ctx, territory, pager.Offset, pager.Limit)

	if err != nil {
	    return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

	for index, value := range tracking {
		profile, err := h.EscortProfileRepository.GetEscortProfile(ctx, value.EscortId)

		if err != nil {
		    continue
        }

		tracking[index].FirstName = profile.FirstName
		tracking[index].LastName = profile.LastName
		tracking[index].Avatar = fmt.Sprintf(
		    "%s/%s/%s",
		    config.InitS3().Endpoint,
		    config.InitS3().Buckets.Profile,
		    profile.Avatar,
		)
	}

	totalRows, _ := h.Repository.CountEscortLocationByTerritory(ctx)
	pagerResult := types.PagerResult{
	    Total: totalRows,
	    Data: tracking,
    }

	return c.JSON(http.StatusOK, pagerResult.GetPagerResult(pager))
}

func (h *EscortTrackingController) GetEscortLocation(c echo.Context) error {
	tracking, err := h.Repository.GetEscortTracking(
	    c.Request().Context(),
	    c.Request().Header.Get("user-id"),
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, tracking)
}

func (h *EscortTrackingController) SetEscortLocation(c echo.Context) (err error) {
	escortTracking := models.EscortTracking{}
	userId := c.Request().Header.Get("user-id")
	escortTracking.EscortId = userId

	if err = c.Bind(&escortTracking.Location); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ctx := c.Request().Context()

	if err = h.Repository.AlterEscortTracking(ctx, &escortTracking); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go h.emitAcknowledge(ctx, userId, &wg)

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
		Operation:   config.InitOperationConfig().NewUser,
		UserId:      userId,
		UserType:    "Escort",
	}
	bytes, _ := json.Marshal(countUserEvent)

	if err := h.KafkaService.SendMessage(config.InitKafkaConfig().Topics.OperationTopic, bytes); err != nil {
		log.Errorf("escort_controller->emitAcknowledge->SendMessage:%s", err.Error())
	}
}
