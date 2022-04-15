package repositories

import (
	"context"
	"escort-book-tracking/db"
	"escort-book-tracking/models"
)

type IEscortProfileRepository interface {
	GetEscortProfile(ctx context.Context, id string) (*models.EscortProfile, error)
}

type EscortProfileRepository struct {
	Data *db.Data
}

func (r *EscortProfileRepository) GetEscortProfile(ctx context.Context, id string) (*models.EscortProfile, error) {
	query := `SELECT a.first_name, a.last_name, b.path
			  FROM profile as a
			  JOIN avatar as b
			  ON a.escort_id = b.escort_id
			  WHERE a.escort_id = $1;`
	row := r.Data.DB.QueryRowContext(ctx, query, id)

	var profile models.EscortProfile
	err := row.Scan(&profile.FirstName, &profile.LastName, &profile.Avatar)

	if err != nil {
		return &models.EscortProfile{}, err
	}

	return &profile, nil
}
