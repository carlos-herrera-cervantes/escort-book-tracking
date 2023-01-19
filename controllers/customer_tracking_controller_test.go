package controllers

import (
    "context"
    "errors"
    "net/http"
    "net/http/httptest"
    "strings"
    "sync"
    "testing"

    "escort-book-tracking/models"
    mockRepositories "escort-book-tracking/repositories/mocks"
    mockServices "escort-book-tracking/services/mocks"

    "github.com/golang/mock/gomock"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
)

func TestCustomerTrackingControllerGetCustomerLocation(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockCustomerTrackingRepository := mockRepositories.NewMockICustomerTrackingRepository(ctrl)
    mockKafkaService := mockServices.NewMockIKafkaService(ctrl)

    customerTrackingController := CustomerTrackingController{
        Repository: mockCustomerTrackingRepository,
        KafkaService: mockKafkaService,
    }

    t.Run("Should return 404 http error", func(t *testing.T) {
        e := echo.New()

        req := httptest.NewRequest(
            http.MethodGet,
            "/api/v1/tracking/customer",
            nil,
        )
        recorder := httptest.NewRecorder()

        c := e.NewContext(req, recorder)

        mockCustomerTrackingRepository.
            EXPECT().
            GetCustomerTracking(gomock.Any(), gomock.Any()).
            Return(&models.CustomerTracking{}, errors.New("dummy error")).
            Times(1)

        res := customerTrackingController.GetCustomerLocation(c)

        assert.Error(t, res)
    })

    t.Run("Should return 200 status code", func(t *testing.T) {
        e := echo.New()

        req := httptest.NewRequest(
            http.MethodGet,
            "/api/v1/tracking/customer",
            nil,
        )
        recorder := httptest.NewRecorder()

        c := e.NewContext(req, recorder)

        mockCustomerTrackingRepository.
            EXPECT().
            GetCustomerTracking(gomock.Any(), gomock.Any()).
            Return(&models.CustomerTracking{}, nil).
            Times(1)

        res := customerTrackingController.GetCustomerLocation(c)

        assert.NoError(t, res)
        assert.Equal(t, http.StatusOK, recorder.Code)
    })
}

func TestCustomerTrackingControllerSetCustomerLocation(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockCustomerTrackingRepository := mockRepositories.NewMockICustomerTrackingRepository(ctrl)
    mockKafkaService := mockServices.NewMockIKafkaService(ctrl)

    customerTrackingController := CustomerTrackingController{
        Repository: mockCustomerTrackingRepository,
        KafkaService: mockKafkaService,
    }

    t.Run("Should return 500 http error", func(t *testing.T) {
        e := echo.New()
        reqBody := `
            {
                "location": {
                    "latitude": 4.6788055,
                    "longitude": -74.0554172
                }
            }
        `

        req := httptest.NewRequest(
            http.MethodPost,
            "/api/v1/tracking/customer",
            strings.NewReader(reqBody),
        )
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        req.Header.Set("user-id", "63915d9be04ff7e3da92ffa3")
        recorder := httptest.NewRecorder()

        c := e.NewContext(req, recorder)

        mockCustomerTrackingRepository.
            EXPECT().
            AlterCustomerTracking(gomock.Any(), gomock.Any()).
            Return(errors.New("dummy error")).
            Times(1)

        res := customerTrackingController.SetCustomerLocation(c)

        assert.Error(t, res)
    })

    t.Run("Should return 201 status code", func(t *testing.T) {
        e := echo.New()
        reqBody := `
        {
            "location": {
                "latitude": 4.6788055,
                "longitude": -74.0554172
            }
        }
        `

        req := httptest.NewRequest(
            http.MethodPost,
            "/api/v1/tracking/customer",
            strings.NewReader(reqBody),
        )
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        req.Header.Set("user-id", "63915d9be04ff7e3da92ffa3")
        recorder := httptest.NewRecorder()

        c := e.NewContext(req, recorder)

        mockCustomerTrackingRepository.
            EXPECT().
            AlterCustomerTracking(gomock.Any(), gomock.Any()).
            Return(nil).
            Times(1)
        mockCustomerTrackingRepository.
            EXPECT().
            GetCustomerTracking(gomock.Any(), gomock.Any()).
            Return(&models.CustomerTracking{}, errors.New("dummy error")).
            Times(1)

        res := customerTrackingController.SetCustomerLocation(c)

        assert.NoError(t, res)
        assert.Equal(t, http.StatusCreated, recorder.Code)
    })
}

func TestCustomerTrackingControllerEmitAcknowledge(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockCustomerTrackingRepository := mockRepositories.NewMockICustomerTrackingRepository(ctrl)
    mockKafkaService := mockServices.NewMockIKafkaService(ctrl)

    customerTrackingController := CustomerTrackingController{
        Repository: mockCustomerTrackingRepository,
        KafkaService: mockKafkaService,
    }

    t.Run("Should exit when acknowledge is true", func(t *testing.T) {
        mockCustomerTrackingRepository.
            EXPECT().
            GetCustomerTracking(gomock.Any(), gomock.Any()).
            Return(&models.CustomerTracking{
                Acknowledged: true,
            }, nil).
            Times(1)
        mockCustomerTrackingRepository.
            EXPECT().
            Acknowledge(gomock.Any(), gomock.Any()).
            Return(nil).
            Times(0)

        var wg sync.WaitGroup
        wg.Add(1)

        ctx := context.Background()
        userId := "63915d9be04ff7e3da92ffa3"
        go customerTrackingController.emitAcknowledge(ctx, userId, &wg)

        wg.Wait()
    })

    t.Run("Should exit when acknowledge fails", func(t *testing.T) {
        mockCustomerTrackingRepository.
            EXPECT().
            GetCustomerTracking(gomock.Any(), gomock.Any()).
            Return(&models.CustomerTracking{
                Acknowledged: false,
            }, nil).
            Times(1)
        mockCustomerTrackingRepository.
            EXPECT().
            Acknowledge(gomock.Any(), gomock.Any()).
            Return(errors.New("dummy error")).
            Times(1)
        mockKafkaService.
            EXPECT().
            SendMessage(gomock.Any(), gomock.Any()).
            Return(nil).
            Times(0)

        var wg sync.WaitGroup
        wg.Add(1)

        ctx := context.Background()
        userId := "63915d9be04ff7e3da92ffa3"
        go customerTrackingController.emitAcknowledge(ctx, userId, &wg)

        wg.Wait()
    })

    t.Run("Should log error when sending message fails", func(t *testing.T) {
        mockCustomerTrackingRepository.
            EXPECT().
            GetCustomerTracking(gomock.Any(), gomock.Any()).
            Return(&models.CustomerTracking{
                Acknowledged: false,
            }, nil).
            Times(1)
        mockCustomerTrackingRepository.
            EXPECT().
            Acknowledge(gomock.Any(), gomock.Any()).
            Return(nil).
            Times(1)
        mockKafkaService.
            EXPECT().
            SendMessage(gomock.Any(), gomock.Any()).
            Return(errors.New("dummy error")).
            Times(1)

        var wg sync.WaitGroup
        wg.Add(1)

        ctx := context.Background()
        userId := "63915d9be04ff7e3da92ffa3"
        go customerTrackingController.emitAcknowledge(ctx, userId, &wg)

        wg.Wait()
    })
}
