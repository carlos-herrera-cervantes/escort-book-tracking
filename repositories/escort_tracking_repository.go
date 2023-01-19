package repositories

import (
	"context"
	"fmt"
	"time"

    "escort-book-tracking/db"
    "escort-book-tracking/models"
)

//go:generate mockgen -destination=./mocks/escort_tracking_repository.go -package=mocks --build_flags=--mod=mod . IEscortTrackingRepository
type IEscortTrackingRepository interface {
	GetEscortTracking(ctx context.Context, id string) (*models.EscortTracking, error)
	GetEscortLocationByTerritory(ctx context.Context, territory string, offset, limit int) ([]models.EscortTracking, error)
	AlterEscortTracking(ctx context.Context, tracking *models.EscortTracking) error
	CountEscortLocationByTerritory(ctx context.Context) (int, error)
	Acknowledge(ctx context.Context, escortId string) error
}

type EscortTrackingRepository struct {
	Data *db.PostgresClient
}

func (r *EscortTrackingRepository) GetEscortTracking(ctx context.Context, id string) (*models.EscortTracking, error) {
	query := `
        SELECT a.id, a.escort_id, st_asgeojson(a.location), a.created_at, a.updated_at, b.name, a.acknowledged
		FROM escort_tracking AS a
		JOIN escort_tracking_status AS b
		ON a.escort_tracking_status_id = b.id
		WHERE escort_id = $1;
    `
	row := r.Data.EscortTrackingDB.QueryRowContext(ctx, query, id)

	var tracking models.EscortTracking
	var stringPoint string

	if err := row.Scan(
		&tracking.Id,
		&tracking.EscortId,
		&stringPoint,
		&tracking.CreatedAt,
		&tracking.UpdatedAt,
		&tracking.EscortTrackingStatus,
		&tracking.Acknowledged,
	); err != nil {
	    return &tracking, err
    }

	tracking.Location.ParseGeoJson(stringPoint)

	return &tracking, nil
}

func (r *EscortTrackingRepository) GetEscortLocationByTerritory(ctx context.Context, territory string, offset, limit int) ([]models.EscortTracking, error) {
	query := `
        SELECT a.id, a.escort_id, st_asgeojson(a.location), a.created_at, a.updated_at, c.name
		FROM escort_tracking AS a
		JOIN territory AS b
		ON st_intersects(a.location, b.location)
		JOIN escort_tracking_status AS c
		ON a.escort_tracking_status_id = c.id
		WHERE b.name = $1 AND c.name IN('Free', 'Busy') OFFSET($2) LIMIT($3);
    `

	rows, err := r.Data.EscortTrackingDB.QueryContext(ctx, query, territory, offset, limit)
	var tracking []models.EscortTracking

	if err != nil {
		return tracking, err
	}

	for rows.Next() {
		var point models.EscortTracking
		var stringPoint string

		if err := rows.Scan(
		    &point.Id,
		    &point.EscortId,
		    &stringPoint,
		    &point.CreatedAt,
		    &point.UpdatedAt,
		    &point.EscortTrackingStatus,
		); err != nil {
		    continue
        }

        point.Location.ParseGeoJson(stringPoint)
		tracking = append(tracking, point)
	}

	return tracking, nil
}

func (r *EscortTrackingRepository) AlterEscortTracking(ctx context.Context, tracking *models.EscortTracking) error {
	query := "SELECT id FROM escort_tracking WHERE escort_id = $1;"
	row := r.Data.EscortTrackingDB.QueryRowContext(ctx, query, tracking.EscortId)

	var id string

	if err := row.Scan(&id); err == nil {
		update := fmt.Sprintf(
			"UPDATE escort_tracking SET location = 'POINT(%f %f)', updated_at = $1 WHERE id = $2;",
			tracking.Location.Latitude,
			tracking.Location.Longitude,
		)

		if _, err = r.Data.EscortTrackingDB.ExecContext(ctx, update, time.Now().UTC(), id); err != nil {
		    return err
        }

		return nil
	}

	point := fmt.Sprintf(`POINT(%f %f)`, tracking.Location.Latitude, tracking.Location.Longitude)
	insert := "INSERT INTO escort_tracking VALUES ($1, $2, $3, $4, $5);"
	tracking.SetDefaultValues()

	if _, err := r.Data.EscortTrackingDB.ExecContext(
		ctx,
		insert,
		tracking.Id,
		tracking.EscortId,
		point,
		time.Now().UTC(),
		time.Now().UTC(),
	); err != nil {
	    return err
    }

	return nil
}

func (r *EscortTrackingRepository) CountEscortLocationByTerritory(ctx context.Context) (int, error) {
	query := `
        SELECT COUNT(a.*) from escort_tracking AS a
		JOIN territory AS b
		ON st_intersects(a.location, b.location)
		WHERE b.name = $1;
    `

	row := r.Data.EscortTrackingDB.QueryRowContext(ctx, query)

	var number int

	if err := row.Scan(&number); err != nil {
	    return 0, err
    }

	return number, nil
}

func (r *EscortTrackingRepository) Acknowledge(ctx context.Context, escortId string) error {
	query := "UPDATE escort_tracking SET acknowledged = TRUE WHERE escort_id = $1;"

	if _, err := r.Data.EscortTrackingDB.ExecContext(ctx, query, escortId); err != nil {
		return err
	}

	return nil
}
