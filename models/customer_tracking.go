package models

import (
	"time"

    "escort-book-tracking/types"

	"github.com/google/uuid"
)

type CustomerTracking struct {
	Id           string         `json:"id"`
	CustomerId   string         `json:"customerId"`
	Location     types.Location `json:"location"`
	Acknowledged bool           `json:"-"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}

func (c *CustomerTracking) SetDefaultValues() { c.Id = uuid.NewString() }
