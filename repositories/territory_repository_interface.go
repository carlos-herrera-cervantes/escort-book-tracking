package repositories

import (
	"context"
	"escort-book-tracking/models"
)

type ITerritoryRepository interface {
	GetOneByName(ctx context.Context, name string) (models.Territory, error)
}
