package repositories

import (
	"context"
	"regexp"
	"testing"
	"time"

    "escort-book-tracking/db"
    "escort-book-tracking/models"
    "escort-book-tracking/types"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetEscortTracking(t *testing.T) {
	database, mock, _ := sqlmock.New()
	data := &db.PostgresClient{
		EscortTrackingDB: database,
	}

	repository := &EscortTrackingRepository{
		Data: data,
	}

	now := time.Now().UTC()
	rows := sqlmock.
		NewRows([]string{"id", "escort_id", "location", "created_at", "updated_at", "name", "acknowledged"}).
		AddRow(
			"dummy-uuid",
			"dummy-escort-id",
			`{"type":"Point","coordinates":[-99.905932819,16.820014824]}`,
			now,
			now,
			"Free",
			false,
		)

	query := `
        SELECT a.id, a.escort_id, st_asgeojson(a.location), a.created_at, a.updated_at, b.name, a.acknowledged
        FROM escort_tracking AS a
        JOIN escort_tracking_status AS b
        ON a.escort_tracking_status_id = b.id
        WHERE escort_id = $1;
    `

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
	ctx := context.Background()
	tracking, err := repository.GetEscortTracking(ctx, "dummy-uuid")

	expectedTracking := &models.EscortTracking{
		Id:       "dummy-uuid",
		EscortId: "dummy-escort-id",
		Location: types.Location{
			Latitude:  -99.905932819,
			Longitude: 16.820014824,
		},
		CreatedAt:            now,
		UpdatedAt:            now,
		EscortTrackingStatus: "Free",
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedTracking, tracking)
}

func TestGetEscortLocationByTerritory(t *testing.T) {
	database, mock, _ := sqlmock.New()
	data := &db.PostgresClient{
		EscortTrackingDB: database,
	}

	repository := &EscortTrackingRepository{
		Data: data,
	}

	now := time.Now().UTC()
	rows := sqlmock.
		NewRows([]string{"id", "escort_id", "location", "created_at", "updated_at", "name"}).
		AddRow(
			"dummy-uuid",
			"dummy-escort-id",
			`{"type":"Point","coordinates":[-99.905932819,16.820014824]}`,
			now,
			now,
			"Free",
		)

	query := `
        SELECT a.id, a.escort_id, st_asgeojson(a.location), a.created_at, a.updated_at, c.name
        FROM escort_tracking AS a
        JOIN territory AS b
        ON st_intersects(a.location, b.location)
        JOIN escort_tracking_status AS c
        ON a.escort_tracking_status_id = c.id
        WHERE b.name = $1 AND c.name IN('Free', 'Busy') OFFSET($2) LIMIT($3);
    `

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
	ctx := context.Background()
	trackings, err := repository.GetEscortLocationByTerritory(ctx, "dummy-territory", 0, 10)

	expectedTrackings := []models.EscortTracking{
		{
			Id:       "dummy-uuid",
			EscortId: "dummy-escort-id",
			Location: types.Location{
				Latitude:  -99.905932819,
				Longitude: 16.820014824,
			},
			CreatedAt:            now,
			UpdatedAt:            now,
			EscortTrackingStatus: "Free",
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedTrackings, trackings)
}
