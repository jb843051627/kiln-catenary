package regression

import (
    "context"
    "errors"
    "testing"
    "time"
    "github.com/jb843051627/kiln-catenary/internal/model"
    "github.com/jb843051627/kiln-catenary/internal/report"
)

func TestBug18_ThermalWindowCloneOwnsSamples(t *testing.T) {
	values := []model.AtmosphereSample{{ID: "late", ObservedAt: time.Unix(2, 0).UTC(), Temperature: 900}, {ID: "early", ObservedAt: time.Unix(1, 0).UTC(), Temperature: 800}}
	window := report.NewThermalWindow("kiln", time.Unix(0, 0).UTC(), time.Unix(3, 0).UTC(), values); ctx, cancel := context.WithCancel(context.Background()); cancel(); _, err := window.CloneContext(ctx); if !errors.Is(err, context.Canceled) { t.Fatalf("canceled thermal clone returned %v", err) }
}
