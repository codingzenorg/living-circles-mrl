package server_test

import (
	"testing"

	"github.com/codingzen/living-circles-mrl/src/server/simulation"
	"github.com/codingzen/living-circles-mrl/src/server/transport"
)

func TestDefaultExpandedSnapshotTransportMeasurementIsDeterministic(t *testing.T) {
	first, err := transport.MeasureSnapshotTransport(simulation.NewSession().Snapshot(), transport.DefaultTickEvery)
	if err != nil {
		t.Fatalf("measure first snapshot transport: %v", err)
	}

	second, err := transport.MeasureSnapshotTransport(simulation.NewSession().Snapshot(), transport.DefaultTickEvery)
	if err != nil {
		t.Fatalf("measure second snapshot transport: %v", err)
	}

	if first != second {
		t.Fatalf("expected deterministic transport measurement, first=%+v second=%+v", first, second)
	}
	if first.PayloadBytes <= 0 {
		t.Fatalf("expected payload bytes to be positive, got %+v", first)
	}
	if first.SnapshotsPerSecond != 10 {
		t.Fatalf("expected snapshots per second 10, got %+v", first)
	}
	if first.ApproxBytesPerSecond != float64(first.PayloadBytes)*first.SnapshotsPerSecond {
		t.Fatalf("expected bytes per second to match payload*rate, got %+v", first)
	}
}

func TestLargerWorldScenarioTransportMeasurementExceedsDefaultExpandedBaseline(t *testing.T) {
	defaultMeasurement, err := transport.MeasureSnapshotTransport(simulation.NewSession().Snapshot(), transport.DefaultTickEvery)
	if err != nil {
		t.Fatalf("measure default snapshot transport: %v", err)
	}

	largerScenario := simulation.NewSessionWithConfig(simulation.Config{
		WorldWidth:                simulation.DefaultExpandedWorldWidth * 2,
		WorldHeight:               simulation.DefaultExpandedWorldHeight * 2,
		UseExpandedPopulation:     true,
		ExpandedAutonomousCount:   simulation.DefaultExpandedAutonomousCount * 2,
		ExpandedFoodCount:         simulation.DefaultExpandedFoodCount * 2,
		PlayerShape:               simulation.DefaultPlayerShape,
		AutonomousShape:           simulation.DefaultPlayerShape,
		SecondaryAutonomousShape:  simulation.DefaultAutoShape,
		PlayerEnergy:              simulation.DefaultPlayerEnergy,
		PlayerChildrenCount:       1,
		AutonomousEnergy:          80,
		SecondaryAutonomousEnergy: simulation.DefaultPlayerEnergy,
		AutonomousChildrenCount:   1,
	})

	largerMeasurement, err := transport.MeasureSnapshotTransport(largerScenario.Snapshot(), transport.DefaultTickEvery)
	if err != nil {
		t.Fatalf("measure larger snapshot transport: %v", err)
	}

	if largerMeasurement.PayloadBytes <= defaultMeasurement.PayloadBytes {
		t.Fatalf("expected larger scenario payload %d to exceed default payload %d", largerMeasurement.PayloadBytes, defaultMeasurement.PayloadBytes)
	}
	if largerMeasurement.ApproxBytesPerSecond <= defaultMeasurement.ApproxBytesPerSecond {
		t.Fatalf("expected larger scenario bytes/sec %v to exceed default bytes/sec %v", largerMeasurement.ApproxBytesPerSecond, defaultMeasurement.ApproxBytesPerSecond)
	}
}
