package repositories

import (
	"context"
	"escort-book-tracking/models"
)

type IEscortTrackingRepository interface {
	GetOne(ctx context.Context, id string) (models.EscortTracking, error)
	GetByTerritory(ctx context.Context, territory string) ([]models.EscortTracking, error)
	Create(ctx context.Context, tracking *models.EscortTracking) error
}
