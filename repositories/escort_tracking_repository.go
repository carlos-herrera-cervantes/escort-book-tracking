package repositories

import (
	"context"
	"escort-book-tracking/db"
	"escort-book-tracking/models"
	"fmt"
	"time"
)

type EscortTrackingRepository struct {
	Data *db.Data
}

func (r *EscortTrackingRepository) GetOne(ctx context.Context, id string) (models.EscortTracking, error) {
	query := `SELECT id, escort_id, st_asgeojson(location), created_at, updated_at
		      FROM escort_tracking WHERE id = $1;`
	row := r.Data.DB.QueryRowContext(ctx, query, id)

	var tracking models.EscortTracking
	var stringPoint string
	err := row.Scan(
		&tracking.Id,
		&tracking.EscortId,
		&stringPoint,
		&tracking.CreatedAt,
		&tracking.UpdatedAt)

	if err != nil {
		return models.EscortTracking{}, err
	}

	tracking.Location.ParseGeoJson(stringPoint)

	return tracking, nil
}

func (r *EscortTrackingRepository) GetByTerritory(ctx context.Context, territory string, offset, limit int) ([]models.EscortTracking, error) {
	query := `SELECT a.* from escort_tracking as a
			  join territory as b
		      on st_intersects(a.location, b.location)
		      where b.name = $1`

	rows, _ := r.Data.DB.QueryContext(ctx, query, territory)
	var trackings []models.EscortTracking

	for rows.Next() {
		var tracking models.EscortTracking
		var stringPoint string

		rows.Scan(&tracking.Id, &tracking.EscortId, &stringPoint, &tracking.CreatedAt, &tracking.UpdatedAt)
		tracking.Location.ParseGeoJson(stringPoint)
		trackings = append(trackings, tracking)
	}

	return trackings, nil
}

func (r *EscortTrackingRepository) UpsertOne(ctx context.Context, tracking *models.EscortTracking) error {
	query := "SELECT id FROM escort_tracking WHERE escort_id = $1;"
	row := r.Data.DB.QueryRowContext(ctx, query, tracking.EscortId)

	var id string
	err := row.Scan(&id)

	if err == nil {
		tracking.Id = id
		update := fmt.Sprintf(
			"UPDATE escort_tracking SET location = 'POINT(%f %f)', updated_at = $1 WHERE id = $2;",
			tracking.Location.Latitude,
			tracking.Location.Longitude,
		)
		r.Data.DB.ExecContext(ctx, update, time.Now().UTC(), id)

		return nil
	}

	point := fmt.Sprintf(`POINT(%f %f)`, tracking.Location.Latitude, tracking.Location.Longitude)
	insert := "INSERT INTO escort_tracking VALUES ($1, $2, $3, $4, $5);"
	tracking.SetDefaultValues()

	r.Data.DB.ExecContext(
		ctx,
		insert,
		tracking.Id,
		tracking.EscortId,
		point,
		time.Now().UTC(),
		time.Now().UTC())

	return nil
}

func (r *EscortTrackingRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(a.*) from escort_tracking as a
			  join territory as b
		      on st_intersects(a.location, b.location)
			  where b.name = $1`

	row := r.Data.DB.QueryRowContext(ctx, query)

	var number int

	row.Scan(&number)

	return number, nil
}
