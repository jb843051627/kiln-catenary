package model

import "strings"

type Kiln struct {
	ID             string
	Code           string
	Cell           string
	MaxTemperature float64
	MaxPressure    float64
	Atmosphere     string
	Active         bool
	CreatedAt      string
}

func (k Kiln) Normalize() Kiln {
	k.Code = strings.ToUpper(strings.TrimSpace(k.Code))
	k.Cell = strings.TrimSpace(k.Cell)
	k.Atmosphere = strings.ToLower(strings.TrimSpace(k.Atmosphere))
	return k
}

func (k Kiln) Valid() bool {
	return k.ID != "" && k.Code != "" && k.Cell != "" && k.MaxTemperature > 0 && k.MaxPressure > 0 && k.Atmosphere != ""
}

func (k Kiln) Accepts(temperature, pressure float64) bool {
	return temperature >= 0 && temperature <= k.MaxTemperature && pressure >= 0 && pressure <= k.MaxPressure
}

func (k Kiln) Safe(temperature, pressure float64) bool {
	return k.Active && k.Accepts(temperature, pressure)
}
