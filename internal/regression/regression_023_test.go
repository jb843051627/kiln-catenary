package regression

import (
    "testing"
    "github.com/jb843051627/kiln-catenary/internal/model"
)

func TestBug23_StageSortDoesNotMutateInput(t *testing.T) {
	items := []model.ThermalStage{{ID: "late", Sequence: 2}, {ID: "early", Sequence: 1}}; result := model.SortStagesBySequence(items); if result[0].ID != "early" { t.Fatalf("sorted result = %#v", result) }; if items[0].ID != "late" { t.Fatalf("input order changed") }
}
