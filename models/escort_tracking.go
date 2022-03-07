package models

import (
	"escort-book-tracking/types"
	"time"

	"github.com/google/uuid"
)

type EscortTracking struct {
	Id        string         `json:"id"`
	EscortId  string         `json:"escortId"`
	Location  types.Location `json:"location"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

func (e *EscortTracking) SetDefaultValues() *EscortTracking {
	e.Id = uuid.NewString()
	return e
}
