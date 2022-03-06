package repositories

import (
	"context"
	"escort-book-tracking/db"
	"escort-book-tracking/models"
	"time"
)

type EscortTrackingRepository struct {
	Data *db.Data
}

func (r *EscortTrackingRepository) GetOne(ctx context.Context, id string) (models.EscortTracking, error) {
	query := "SELECT * FROM escorttracking WHERE id = $1;"
	row := r.Data.DB.QueryRowContext(ctx, query, id)

	var tracking models.EscortTracking
	err := row.Scan(
		&tracking.Id,
		&tracking.EscortId,
		&tracking.Location,
		&tracking.CreatedAt,
		&tracking.UpdatedAt)

	if err != nil {
		return models.EscortTracking{}, err
	}

	return tracking, nil
}

func (r *EscortTrackingRepository) GetByTerritory(ctx context.Context, territory string) ([]models.EscortTracking, error) {
	query := `SELECT a.* from customer_tracking as a
			  join territory as b
		      on st_intersects(a.location, b.location)
		      where b.name = $1`

	rows, _ := r.Data.DB.QueryContext(ctx, query, territory)
	var trackings []models.EscortTracking

	for rows.Next() {
		var tracking models.EscortTracking

		rows.Scan(&tracking.Id, &tracking.EscortId, &tracking.Location, &tracking.CreatedAt, &tracking.UpdatedAt)
		trackings = append(trackings, tracking)
	}

	return trackings, nil
}

func (r *EscortTrackingRepository) Create(ctx context.Context, tracking *models.EscortTracking) error {
	query := "INSERT INTO escorttracking VALUES ($1, $2, $3, $4, $5);"
	tracking.SetDefaultValues()

	_, err := r.Data.DB.ExecContext(
		ctx,
		query,
		tracking.Id,
		tracking.EscortId,
		tracking.Location,
		time.Now().UTC(),
		time.Now().UTC())

	if err != nil {
		return err
	}

	return nil
}
