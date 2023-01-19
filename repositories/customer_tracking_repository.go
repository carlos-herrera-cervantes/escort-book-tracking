package repositories

import (
	"context"
	"fmt"
	"time"

    "escort-book-tracking/db"
    "escort-book-tracking/models"
)

//go:generate mockgen -destination=./mocks/customer_tracking_repository.go -package=mocks --build_flags=--mod=mod . ICustomerTrackingRepository
type ICustomerTrackingRepository interface {
	GetCustomerTracking(ctx context.Context, id string) (*models.CustomerTracking, error)
	AlterCustomerTracking(ctx context.Context, tracking *models.CustomerTracking) error
	Acknowledge(ctx context.Context, customerId string) error
}

type CustomerTrackingRepository struct {
	Data *db.PostgresClient
}

func (r *CustomerTrackingRepository) GetCustomerTracking(ctx context.Context, id string) (*models.CustomerTracking, error) {
	query := `
        SELECT id, customer_id, st_asgeojson(location), created_at, updated_at, acknowledged
	    FROM customer_tracking WHERE customer_id = $1;
    `
	row := r.Data.EscortTrackingDB.QueryRowContext(ctx, query, id)

	var tracking models.CustomerTracking
	var stringPoint string

	if err := row.Scan(
		&tracking.Id,
		&tracking.CustomerId,
		&stringPoint,
		&tracking.CreatedAt,
		&tracking.UpdatedAt,
		&tracking.Acknowledged,
	); err != nil {
		return &tracking, err
	}

	tracking.Location.ParseGeoJson(stringPoint)

	return &tracking, nil
}

func (r *CustomerTrackingRepository) AlterCustomerTracking(ctx context.Context, tracking *models.CustomerTracking) error {
	query := "SELECT id FROM customer_tracking WHERE customer_id = $1;"
	row := r.Data.EscortTrackingDB.QueryRowContext(ctx, query, tracking.CustomerId)

	var id string

	if err := row.Scan(&id); err == nil {
		update := fmt.Sprintf(
			"UPDATE customer_tracking SET location = 'POINT(%f %f)', updated_at = $1 WHERE id = $2;",
			tracking.Location.Latitude,
			tracking.Location.Longitude,
		)

		if _, err := r.Data.EscortTrackingDB.ExecContext(ctx, update, time.Now().UTC(), id); err != nil {
		    return err
        }

		return nil
	}

	point := fmt.Sprintf(`POINT(%f %f)`, tracking.Location.Latitude, tracking.Location.Longitude)
	insert := "INSERT INTO customer_tracking VALUES ($1, $2, $3, $4, $5);"
	tracking.SetDefaultValues()

	if _, err := r.Data.EscortTrackingDB.ExecContext(
		ctx,
		insert,
		tracking.Id,
		tracking.CustomerId,
		point,
		time.Now().UTC(),
		time.Now().UTC(),
	); err != nil {
	    return err
    }

	return nil
}

func (r *CustomerTrackingRepository) Acknowledge(ctx context.Context, customerId string) error {
	query := "UPDATE customer_tracking SET acknowledged = TRUE WHERE customer_id = $1;"

	if _, err := r.Data.EscortTrackingDB.ExecContext(ctx, query, customerId); err != nil {
		return err
	}

	return nil
}
