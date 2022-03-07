package repositories

import (
	"context"
	"escort-book-tracking/models"
)

type IEscortTrackingRepository interface {
	GetOne(ctx context.Context, id string) (models.EscortTracking, error)
	GetByTerritory(ctx context.Context, territory string, offset, limit int) ([]models.EscortTracking, error)
	UpsertOne(ctx context.Context, tracking *models.EscortTracking) error
	Count(ctx context.Context) (int, error)
}
