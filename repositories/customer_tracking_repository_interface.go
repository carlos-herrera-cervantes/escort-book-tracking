package repositories

import (
	"context"
	"escort-book-tracking/models"
)

type ICustomerTrackingRepository interface {
	GetOne(ctx context.Context, id string) (models.CustomerTracking, error)
	Create(ctx context.Context, tracking *models.CustomerTracking) error
}
