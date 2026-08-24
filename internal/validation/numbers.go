package validation

import (
	"fmt"
	"math"
)

func Temperature(value, max float64) error {
	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || value > max {
		return fmt.Errorf("temperature is outside kiln range")
	}
	return nil
}

func Pressure(value, max float64) error {
	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || value > max {
		return fmt.Errorf("pressure is outside kiln range")
	}
	return nil
}

func Positive(value int, name string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive", name)
	}
	return nil
}
