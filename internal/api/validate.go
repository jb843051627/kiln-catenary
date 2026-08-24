package api

import (
	"fmt"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"github.com/jb843051627/kiln-catenary/internal/validation"
	"strings"
	"time"
)

func KilnModel(in KilnRequest) (model.Kiln, error) {
	in.Code = strings.TrimSpace(in.Code)
	in.Cell = strings.TrimSpace(in.Cell)
	in.Atmosphere = strings.TrimSpace(in.Atmosphere)
	if in.Code == "" || in.Cell == "" || in.Atmosphere == "" {
		return model.Kiln{}, fmt.Errorf("code, cell and atmosphere are required")
	}
	if in.MaxTemperature <= 0 || in.MaxPressure <= 0 {
		return model.Kiln{}, fmt.Errorf("kiln limits must be positive")
	}
	return model.Kiln{Code: in.Code, Cell: in.Cell, MaxTemperature: in.MaxTemperature, MaxPressure: in.MaxPressure, Atmosphere: in.Atmosphere, Active: true}, nil
}
func ZoneModel(in ZoneRequest) (model.Zone, error) {
	if in.KilnID == "" || in.Name == "" || in.Target < 0 || in.RampRate <= 0 || in.DeviationLimit <= 0 {
		return model.Zone{}, fmt.Errorf("zone values are invalid")
	}
	return model.Zone{KilnID: in.KilnID, Name: in.Name, Target: in.Target, RampRate: in.RampRate, DeviationLimit: in.DeviationLimit, Enabled: true}, nil
}
func StageModel(in StageRequest) (model.ThermalStage, error) {
	if err := validation.Positive(in.Sequence, "sequence"); err != nil {
		return model.ThermalStage{}, err
	}
	if in.KilnID == "" || in.Name == "" {
		return model.ThermalStage{}, fmt.Errorf("kiln and stage name are required")
	}
	if in.Kind != model.StageRamp && in.Kind != model.StageHold && in.Kind != model.StageCool {
		return model.ThermalStage{}, fmt.Errorf("stage kind is invalid")
	}
	return model.ThermalStage{KilnID: in.KilnID, Name: in.Name, Kind: in.Kind, Sequence: in.Sequence, StartTemp: in.StartTemp, EndTemp: in.EndTemp, Hold: time.Duration(in.HoldSeconds) * time.Second, Status: model.StageReady, Interlock: in.Interlock}, nil
}
func SampleModel(in SampleRequest) (model.AtmosphereSample, error) {
	at, err := Time(in.ObservedAt)
	if err != nil {
		return model.AtmosphereSample{}, err
	}
	if in.KilnID == "" || in.RunID == "" || in.Sequence < 0 {
		return model.AtmosphereSample{}, fmt.Errorf("sample identity is invalid")
	}
	return model.AtmosphereSample{KilnID: in.KilnID, RunID: in.RunID, Temperature: in.Temperature, Pressure: in.Pressure, Oxidation: in.Oxidation, ObservedAt: at, Sequence: in.Sequence, Source: strings.TrimSpace(in.Source)}, nil
}
