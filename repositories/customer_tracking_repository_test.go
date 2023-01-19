package repositories

import (
    "context"
    "fmt"
    "regexp"
    "testing"
    "time"

    "escort-book-tracking/db"
    "escort-book-tracking/models"
    "escort-book-tracking/types"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
)

func TestGetCustomerTracking(t *testing.T) {
	database, mock, _ := sqlmock.New()
	data := &db.PostgresClient{
		EscortTrackingDB: database,
	}
	repository := &CustomerTrackingRepository{
		Data: data,
	}

	now := time.Now().UTC()
	rows := sqlmock.
		NewRows([]string{"id", "customer_id", "location", "created_at", "updated_at", "acknowledged"}).
		AddRow(
		    "dummy-uuid",
		    "dummy-customer-id",
			`{"type":"Point","coordinates":[-99.905932819,16.820014824]}`,
			now,
			now,
			false,
		)

	query := `
        SELECT id, customer_id, st_asgeojson(location), created_at, updated_at, acknowledged
        FROM customer_tracking WHERE customer_id = $1;
    `

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

func TestAlterCustomerTracking(t *testing.T) {
	database, mock, _ := sqlmock.New()
	data := &db.PostgresClient{
		EscortTrackingDB: database,
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

	query2 := fmt.Sprintf(
        "UPDATE customer_tracking SET location = 'POINT(%f %f)', updated_at = $1 WHERE id = $2;",
        -99.905933,
        16.820015,
    )

	mock.ExpectQuery(regexp.QuoteMeta(query1)).WillReturnRows(rows1)
    mock.ExpectExec(regexp.QuoteMeta(query2)).WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err := repository.AlterCustomerTracking(ctx, inputCustomerTracking)

	assert.Nil(t, err)
}
