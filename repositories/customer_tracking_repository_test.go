package repositories

import (
	"context"
	"escort-book-tracking/db"
	"escort-book-tracking/models"
	"escort-book-tracking/types"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetCustomerTracking(t *testing.T) {
	database, mock, _ := sqlmock.New()
	data := &db.Data{
		DB: database,
	}

	repository := &CustomerTrackingRepository{
		Data: data,
	}

	now := time.Now().UTC()
	rows := sqlmock.
		NewRows([]string{"id", "customer_id", "location", "created_at", "updated_at"}).
		AddRow(
			"dummy-uuid",
			"dummy-customer-id",
			`{"type":"Point","coordinates":[-99.905932819,16.820014824]}`,
			now,
			now,
		)

	query := `SELECT id, customer_id, st_asgeojson(location), created_at, updated_at
			  FROM customer_tracking WHERE customer_id = $1;`

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
	ctx := context.Background()
	tracking, err := repository.GetCustomerTracking(ctx, "dummy-uuid")

	expectedTracking := &models.CustomerTracking{
		Id:         "dummy-uuid",
		CustomerId: "dummy-customer-id",
		Location: types.Location{
			Latitude:  -99.905932819,
			Longitude: 16.820014824,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedTracking, tracking)
}

func TestUpsertCustomerTracking(t *testing.T) {
	database, mock, _ := sqlmock.New()
	data := &db.Data{
		DB: database,
	}

	repository := &CustomerTrackingRepository{
		Data: data,
	}

	inputCustomerTracking := &models.CustomerTracking{
		Location: types.Location{
			Latitude:  -99.905932819,
			Longitude: 16.820014824,
		},
		CustomerId: "dummy-customer-id",
	}

	query1 := "SELECT id FROM customer_tracking WHERE customer_id = $1;"
	rows1 := sqlmock.NewRows([]string{"id"}).AddRow("dummy-uuid")

	mock.ExpectQuery(regexp.QuoteMeta(query1)).WillReturnRows(rows1)

	ctx := context.Background()
	err := repository.UpsertCustomerTracking(ctx, inputCustomerTracking)

	assert.Nil(t, err)
}
