package models

import (
	"escort-book-tracking/types"
	"time"

	"github.com/google/uuid"
)

type EscortTracking struct {
	Id                   string         `json:"id"`
	EscortId             string         `json:"escortId"`
	Location             types.Location `json:"location"`
	EscortTrackingStatus string         `json:"escortTrackingStatus"`
	FirstName            string         `json:"firstName"`
	LastName             string         `json:"lastName"`
	Avatar               string         `json:"avatar"`
	Acknowledged         bool           `json:"-"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
}

func (e *EscortTracking) SetDefaultValues() { e.Id = uuid.NewString() }
