package repositories

import (
	"context"
	"escort-book-tracking/db"
	"escort-book-tracking/models"
	"fmt"
	"time"
)

type ICustomerTrackingRepository interface {
	GetCustomerTracking(ctx context.Context, id string) (*models.CustomerTracking, error)
	UpsertCustomerTracking(ctx context.Context, tracking *models.CustomerTracking) error
	Acknowledge(ctx context.Context, customerId string) error
}

type CustomerTrackingRepository struct {
	Data *db.Data
}

func (r *CustomerTrackingRepository) GetCustomerTracking(ctx context.Context, id string) (*models.CustomerTracking, error) {
	query := `SELECT id, customer_id, st_asgeojson(location), created_at, updated_at, acknowledged
			  FROM customer_tracking WHERE customer_id = $1;`
	row := r.Data.DB.QueryRowContext(ctx, query, id)

	var tracking models.CustomerTracking
	var stringPoint string
	err := row.Scan(
		&tracking.Id,
		&tracking.CustomerId,
		&stringPoint,
		&tracking.CreatedAt,
		&tracking.UpdatedAt,
		&tracking.Acknowledged,
	)

	if err != nil {
		return &models.CustomerTracking{}, err
	}

	tracking.Location.ParseGeoJson(stringPoint)

	return &tracking, nil
}

func (r *CustomerTrackingRepository) UpsertCustomerTracking(ctx context.Context, tracking *models.CustomerTracking) error {
	query := "SELECT id FROM customer_tracking WHERE customer_id = $1;"
	row := r.Data.DB.QueryRowContext(ctx, query, tracking.CustomerId)

	var id string
	err := row.Scan(&id)

	if err == nil {
		tracking.Id = id
		update := fmt.Sprintf(
			"UPDATE customer_tracking SET location = 'POINT(%f %f)', updated_at = $1 WHERE id = $2;",
			tracking.Location.Latitude,
			tracking.Location.Longitude,
		)
		r.Data.DB.ExecContext(ctx, update, time.Now().UTC(), id)

		return nil
	}

	point := fmt.Sprintf(`POINT(%f %f)`, tracking.Location.Latitude, tracking.Location.Longitude)
	insert := "INSERT INTO customer_tracking VALUES ($1, $2, $3, $4, $5);"
	tracking.SetDefaultValues()

	r.Data.DB.ExecContext(
		ctx,
		insert,
		tracking.Id,
		tracking.CustomerId,
		point,
		time.Now().UTC(),
		time.Now().UTC())

	return nil
}

func (r *CustomerTrackingRepository) Acknowledge(ctx context.Context, customerId string) error {
	query := "UPDATE customer_tracking SET acknowledged = TRUE WHERE customer_id = $1;"

	if _, err := r.Data.DB.ExecContext(ctx, query, customerId); err != nil {
		return err
	}

	return nil
}
