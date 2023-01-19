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

func TestGetEscortProfile(t *testing.T) {
	database, mock, _ := sqlmock.New()
	data := &db.PostgresClient{
		EscortProfileDB: database,
	}
	repository := &EscortProfileRepository{
		Data: data,
	}

	rows := sqlmock.
		NewRows([]string{"first_name", "last_name", "avatar"}).
		AddRow("María", "Cruz", "dummy-id/avatar.png")

	query := `
        SELECT a.first_name, a.last_name, b.path
        FROM profile as a
        JOIN avatar as b
        ON a.escort_id = b.escort_id
        WHERE a.escort_id = $1;
    `

	ctx := context.Background()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
	profile, err := repository.GetEscortProfile(ctx, "dummy-id")

	expectedProfile := &models.EscortProfile{
		FirstName: "María",
		LastName:  "Cruz",
		Avatar:    "dummy-id/avatar.png",
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedProfile, profile)
}
