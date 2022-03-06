package repositories

import (
	"context"
	"escort-book-tracking/db"
	"escort-book-tracking/models"
	"time"
)

type CustomerTrackingRepository struct {
	Data *db.Data
}

func (r *CustomerTrackingRepository) GetOne(ctx context.Context, id string) (models.CustomerTracking, error) {
	query := "SELECT * FROM customertracking WHERE id = $1;"
	row := r.Data.DB.QueryRowContext(ctx, query, id)

	var tracking models.CustomerTracking
	err := row.Scan(
		&tracking.Id,
		&tracking.CustomerId,
		&tracking.Location,
		&tracking.CreatedAt,
		&tracking.UpdatedAt)

	if err != nil {
		return models.CustomerTracking{}, err
	}

	return tracking, nil
}

func (r *CustomerTrackingRepository) Create(ctx context.Context, tracking *models.CustomerTracking) error {
	query := "INSERT INTO customertracking VALUES ($1, $2, $3, $4, $5);"
	tracking.SetDefaultValues()

	_, err := r.Data.DB.ExecContext(
		ctx,
		query,
		tracking.Id,
		tracking.CustomerId,
		tracking.Location,
		time.Now().UTC(),
		time.Now().UTC())

	if err != nil {
		return err
	}

	return nil
}
