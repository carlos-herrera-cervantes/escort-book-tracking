package repositories

import (
	"context"

	"escort-book-tracking/db"
	"escort-book-tracking/models"
)

type ITerritoryRepository interface {
	GetTerritoryByName(ctx context.Context, name string) (*models.Territory, error)
}

type TerritoryRepository struct {
	Data *db.PostgresClient
}

func (r *TerritoryRepository) GetTerritoryByName(ctx context.Context, name string) (*models.Territory, error) {
	query := "SELECT id, name FROM territory WHERE name = $1;"
	row := r.Data.EscortTrackingDB.QueryRowContext(ctx, query, name)

	var territory models.Territory

	if err := row.Scan(&territory.Id, &territory.Name); err != nil {
	    return &territory, err
    }

	return &territory, nil
}
