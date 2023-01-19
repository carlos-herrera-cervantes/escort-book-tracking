package controllers

import (
    "context"
    "errors"
    "net/http"
    "net/http/httptest"
    "net/url"
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

func TestEscortTrackingControllerGetLocationsByTerritory(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockEscortTrackingRepository := mockRepositories.NewMockIEscortTrackingRepository(ctrl)
    mockEscortProfileRepository := mockRepositories.NewMockIEscortProfileRepository(ctrl)
    mockKafkaService := mockServices.NewMockIKafkaService(ctrl)

    escortTrackingController := EscortTrackingController{
        Repository: mockEscortTrackingRepository,
        EscortProfileRepository: mockEscortProfileRepository,
        KafkaService: mockKafkaService,
    }

    t.Run("Should return 400 when pager is invalid", func(t *testing.T) {
        e := echo.New()

        query := make(url.Values)
        query.Set("offset", "-1")
        query.Set("limit", "12")

        req := httptest.NewRequest(
            http.MethodGet,
            "/api/v1/tracking/escorts?"+query.Encode(),
            nil,
        )
        recorder := httptest.NewRecorder()

        c := e.NewContext(req, recorder)
        res := escortTrackingController.GetLocationsByTerritory(c)

        assert.Error(t, res)
    })

    t.Run("Should return 500 when DB query fails", func(t *testing.T) {
        e := echo.New()

        query := make(url.Values)
        query.Set("offset", "1")
        query.Set("limit", "10")
        query.Set("territory", "acapulco")

        req := httptest.NewRequest(
            http.MethodGet,
            "/api/v1/tracking/escorts?"+query.Encode(),
            nil,
        )
        recorder := httptest.NewRecorder()

        c := e.NewContext(req, recorder)

        mockEscortTrackingRepository.
            EXPECT().
            GetEscortLocationByTerritory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
            Return([]models.EscortTracking{}, errors.New("dummy error")).
            Times(1)

        res := escortTrackingController.GetLocationsByTerritory(c)

        assert.Error(t, res)
    })

    t.Run("Should return 200 status code", func(t *testing.T) {
        e := echo.New()

        query := make(url.Values)
        query.Set("offset", "1")
        query.Set("limit", "10")
        query.Set("territory", "acapulco")

        req := httptest.NewRequest(
            http.MethodGet,
            "/api/v1/tracking/escorts?"+query.Encode(),
            nil,
        )
        recorder := httptest.NewRecorder()

        c := e.NewContext(req, recorder)

        mockEscortTrackingRepository.
            EXPECT().
            GetEscortLocationByTerritory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
            Return([]models.EscortTracking{
                {
                    EscortId: "63c63a46380e07938a6652c9",
                },
            }, nil).
            Times(1)
        mockEscortProfileRepository.
            EXPECT().
            GetEscortProfile(gomock.Any(), gomock.Any()).
            Return(&models.EscortProfile{
                FirstName: "Test",
                LastName: "Escort",
                Avatar: "/63c63a46380e07938a6652c9/profile.png",
            }, nil).
            Times(1)
        mockEscortTrackingRepository.
            EXPECT().
            CountEscortLocationByTerritory(gomock.Any()).
            Return(1, nil).
            Times(1)

        res := escortTrackingController.GetLocationsByTerritory(c)

        assert.NoError(t, res)
        assert.Equal(t, http.StatusOK, recorder.Code)
    })
}

func TestEscortTrackingControllerGetEscortLocation(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockEscortTrackingRepository := mockRepositories.NewMockIEscortTrackingRepository(ctrl)
    mockEscortProfileRepository := mockRepositories.NewMockIEscortProfileRepository(ctrl)
    mockKafkaService := mockServices.NewMockIKafkaService(ctrl)

    escortTrackingController := EscortTrackingController{
        Repository: mockEscortTrackingRepository,
        EscortProfileRepository: mockEscortProfileRepository,
        KafkaService: mockKafkaService,
    }

    t.Run("Should return 404 http error", func(t *testing.T) {
        e := echo.New()

        req := httptest.NewRequest(
            http.MethodGet,
            "/api/v1/tracking/escort",
            nil,
        )
        req.Header.Set("user-id", "63c63a46380e07938a6652c9")
        recorder := httptest.NewRecorder()

        c := e.NewContext(req, recorder)

        mockEscortTrackingRepository.
            EXPECT().
            GetEscortTracking(gomock.Any(), gomock.Any()).
            Return(&models.EscortTracking{}, errors.New("dummy error")).
            Times(1)

        res := escortTrackingController.GetEscortLocation(c)

        assert.Error(t, res)
    })

    t.Run("Should return 200 status code", func(t *testing.T) {
        e := echo.New()

        req := httptest.NewRequest(
            http.MethodGet,
            "/api/v1/tracking/escort",
            nil,
        )
        req.Header.Set("user-id", "63c63a46380e07938a6652c9")
        recorder := httptest.NewRecorder()

        c := e.NewContext(req, recorder)

        mockEscortTrackingRepository.
            EXPECT().
            GetEscortTracking(gomock.Any(), gomock.Any()).
            Return(&models.EscortTracking{}, nil).
            Times(1)

        res := escortTrackingController.GetEscortLocation(c)

        assert.NoError(t, res)
        assert.Equal(t, http.StatusOK, recorder.Code)
    })
}

func TestEscortTrackingControllerSetEscortLocation(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockEscortTrackingRepository := mockRepositories.NewMockIEscortTrackingRepository(ctrl)
    mockEscortProfileRepository := mockRepositories.NewMockIEscortProfileRepository(ctrl)
    mockKafkaService := mockServices.NewMockIKafkaService(ctrl)

    escortTrackingController := EscortTrackingController{
        Repository: mockEscortTrackingRepository,
        EscortProfileRepository: mockEscortProfileRepository,
        KafkaService: mockKafkaService,
    }

    t.Run("Should return 500 http error when alter fails", func(t *testing.T) {
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
            "/api/v1/tracking/escort",
            strings.NewReader(reqBody),
        )
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        req.Header.Set("user-id", "63c63a46380e07938a6652c9")
        recorder := httptest.NewRecorder()

        c := e.NewContext(req, recorder)

        mockEscortTrackingRepository.
            EXPECT().
            AlterEscortTracking(gomock.Any(), gomock.Any()).
            Return(errors.New("dummy error")).
            Times(1)

        res := escortTrackingController.SetEscortLocation(c)

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
            "/api/v1/tracking/escort",
            strings.NewReader(reqBody),
        )
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        req.Header.Set("user-id", "63c63a46380e07938a6652c9")
        recorder := httptest.NewRecorder()

        c := e.NewContext(req, recorder)

        mockEscortTrackingRepository.
            EXPECT().
            AlterEscortTracking(gomock.Any(), gomock.Any()).
            Return(nil).
            Times(1)
        mockEscortTrackingRepository.
            EXPECT().
            GetEscortTracking(gomock.Any(), gomock.Any()).
            Return(&models.EscortTracking{}, errors.New("dummy error")).
            Times(1)

        res := escortTrackingController.SetEscortLocation(c)

        assert.NoError(t, res)
        assert.Equal(t, http.StatusCreated, recorder.Code)
    })
}

func TestEscortTrackingControllerEmitAcknowledge(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockEscortTrackingRepository := mockRepositories.NewMockIEscortTrackingRepository(ctrl)
    mockEscortProfileRepository := mockRepositories.NewMockIEscortProfileRepository(ctrl)
    mockKafkaService := mockServices.NewMockIKafkaService(ctrl)

    escortTrackingController := EscortTrackingController{
        Repository: mockEscortTrackingRepository,
        EscortProfileRepository: mockEscortProfileRepository,
        KafkaService: mockKafkaService,
    }

    t.Run("Should exit when acknowledge is true", func(t *testing.T) {
        mockEscortTrackingRepository.
            EXPECT().
            GetEscortTracking(gomock.Any(), gomock.Any()).
            Return(&models.EscortTracking{
                Acknowledged: true,
            }, nil).
            Times(1)
        mockEscortTrackingRepository.
            EXPECT().
            Acknowledge(gomock.Any(), gomock.Any()).
            Return(nil).
            Times(0)

        var wg sync.WaitGroup
        wg.Add(1)

        ctx := context.Background()
        userId := "63915d9be04ff7e3da92ffa3"
        go escortTrackingController.emitAcknowledge(ctx, userId, &wg)

        wg.Wait()
    })

    t.Run("Should exit when acknowledge fails", func(t *testing.T) {
        mockEscortTrackingRepository.
            EXPECT().
            GetEscortTracking(gomock.Any(), gomock.Any()).
            Return(&models.EscortTracking{
                Acknowledged: false,
            }, nil).
            Times(1)
        mockEscortTrackingRepository.
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
        go escortTrackingController.emitAcknowledge(ctx, userId, &wg)

        wg.Wait()
    })

    t.Run("Should log error when sending message fails", func(t *testing.T) {
        mockEscortTrackingRepository.
            EXPECT().
            GetEscortTracking(gomock.Any(), gomock.Any()).
            Return(&models.EscortTracking{
                Acknowledged: false,
            }, nil).
            Times(1)
        mockEscortTrackingRepository.
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
        go escortTrackingController.emitAcknowledge(ctx, userId, &wg)

        wg.Wait()
    })
}
