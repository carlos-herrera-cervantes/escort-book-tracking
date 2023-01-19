package repositories

import (
	"context"
	"regexp"
	"testing"

    "escort-book-tracking/db"
    "escort-book-tracking/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetTerritoryByName(t *testing.T) {
	database, mock, _ := sqlmock.New()
	data := &db.PostgresClient{
		EscortTrackingDB: database,
	}
	repository := &TerritoryRepository{
		Data: data,
	}

	rows := sqlmock.
        NewRows([]string{"id", "name"}).
		AddRow("dummy-uuid", "dummy-territory")

	query := "SELECT id, name FROM territory WHERE name = $1;"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
	ctx := context.Background()

	territory, err := repository.GetTerritoryByName(ctx, "dummy-territory")
	expectedTerritory := &models.Territory{
		Id:   "dummy-uuid",
		Name: "dummy-territory",
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedTerritory, territory)
}
