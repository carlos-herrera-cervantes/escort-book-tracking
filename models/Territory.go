package models

import (
	"time"

	"github.com/google/uuid"
)

type Territory struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	StateId   string    `json:"stateId"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (t *Territory) SetDefaultValues() { t.Id = uuid.NewString() }
