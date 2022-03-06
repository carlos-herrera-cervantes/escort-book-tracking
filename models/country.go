package models

import (
	"time"

	"github.com/google/uuid"
)

type Country struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (c *Country) SetDefaultValues() *Country {
	c.Id = uuid.NewString()
	return c
}
