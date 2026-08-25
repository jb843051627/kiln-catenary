package regression

import (
    "testing"
    "time"
    "github.com/jb843051627/kiln-catenary/internal/model"
    "github.com/jb843051627/kiln-catenary/internal/report"
)

func TestBug22_ThermalWindowOwnsInput(t *testing.T) {
	start := time.Unix(10, 0).UTC(); values := []model.AtmosphereSample{{ID: "late", ObservedAt: start.Add(time.Hour)}, {ID: "early", ObservedAt: start}}; window := report.NewThermalWindow("kiln", start, start.Add(2*time.Hour), values); if window.Samples[0].ID != "early" { t.Fatalf("window order = %#v", window.Samples) }; if values[0].ID != "late" { t.Fatalf("input order changed to %s", values[0].ID) }
}
