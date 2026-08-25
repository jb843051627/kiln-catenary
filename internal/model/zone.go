package model

type Zone struct {
	ID             string
	KilnID         string
	Name           string
	Target         float64
	RampRate       float64
	DeviationLimit float64
	Enabled        bool
}

func (z Zone) Valid() bool {
	return z.ID != "" && z.KilnID != "" && z.Name != "" && z.Target >= 0 && z.RampRate > 0 && z.DeviationLimit > 0
}

func (z Zone) Within(value float64) bool {
	return value >= z.Target-z.DeviationLimit && value <= z.Target+z.DeviationLimit
}

func (z Zone) Delta(value float64) float64 {
	return value - z.Target
}
