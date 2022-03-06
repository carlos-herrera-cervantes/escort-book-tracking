package models

import (
	"time"

	"github.com/google/uuid"
)

type CustomerTracking struct {
	Id         string    `json:"id"`
	CustomerId string    `json:"customerId"`
	Location   string    `json:"location"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (c *CustomerTracking) SetDefaultValues() *CustomerTracking {
	c.Id = uuid.NewString()
	return c
}
