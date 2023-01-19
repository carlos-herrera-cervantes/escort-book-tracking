package repositories

import (
	"context"

	"escort-book-tracking/db"
	"escort-book-tracking/models"
)

//go:generate mockgen -destination=./mocks/escort_profile_repository.go -package=mocks --build_flags=--mod=mod . IEscortProfileRepository
type IEscortProfileRepository interface {
	GetEscortProfile(ctx context.Context, id string) (*models.EscortProfile, error)
}

type EscortProfileRepository struct {
	Data *db.PostgresClient
}

func (r *EscortProfileRepository) GetEscortProfile(ctx context.Context, id string) (*models.EscortProfile, error) {
	query := `
        SELECT a.first_name, a.last_name, b.path
		FROM profile as a
		JOIN avatar as b
		ON a.escort_id = b.escort_id
		WHERE a.escort_id = $1;
    `
	row := r.Data.EscortProfileDB.QueryRowContext(ctx, query, id)

	var profile models.EscortProfile

	if err := row.Scan(&profile.FirstName, &profile.LastName, &profile.Avatar); err != nil {
        return &profile, err
    }

	return &profile, nil
}
