package types

import (
	"encoding/json"
)

type Location struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

func (l *Location) ParseGeoJson(geoJson string) {
	var point map[string]interface{}
	json.Unmarshal([]byte(geoJson), &point)
	coordinates := point["coordinates"].([]interface{})

	l.Latitude = coordinates[0].(float64)
	l.Longitude = coordinates[1].(float64)
}
