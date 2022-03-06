package models

import (
	"time"

	"github.com/google/uuid"
)

type State struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CountryId string    `json:"countryId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s *State) SetDefaultValues() *State {
	s.Id = uuid.NewString()
	return s
}
