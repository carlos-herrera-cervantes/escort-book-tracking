package repositories

import (
	"context"
	"escort-book-tracking/db"
	"escort-book-tracking/models"
)

type TerritoryRepository struct {
	Data *db.Data
}

func (r *TerritoryRepository) GetOneByName(ctx context.Context, name string) (models.Territory, error) {
	query := "SELECT name FROM territory WHERE name = $1;"
	row := r.Data.DB.QueryRowContext(ctx, query, name)

	var territory models.Territory
	err := row.Scan(
		&territory)

	if err != nil {
		return models.Territory{}, err
	}

	return territory, nil
}
