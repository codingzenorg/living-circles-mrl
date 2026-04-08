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

func TestViewportTransportSnapshotKeepsMinimapOrientationWhileCullingLocalDetail(t *testing.T) {
	fullSnapshot := simulation.NewSession().Snapshot()
	transportSnapshot := transport.BuildViewportSnapshot(fullSnapshot)

	if transportSnapshot.Player == nil {
		t.Fatal("expected player to remain present in transport snapshot")
	}
	if len(transportSnapshot.AutonomousCircles) >= len(fullSnapshot.AutonomousCircles) {
		t.Fatalf("expected local autonomous detail to be culled below full world count, local=%d full=%d", len(transportSnapshot.AutonomousCircles), len(fullSnapshot.AutonomousCircles))
	}
	if len(transportSnapshot.Foods) >= len(fullSnapshot.Foods) {
		t.Fatalf("expected local food detail to be culled below full world count, local=%d full=%d", len(transportSnapshot.Foods), len(fullSnapshot.Foods))
	}
	if transportSnapshot.TotalAutonomousCircles != len(fullSnapshot.AutonomousCircles) {
		t.Fatalf("expected total autonomous count %d, got %d", len(fullSnapshot.AutonomousCircles), transportSnapshot.TotalAutonomousCircles)
	}
	if transportSnapshot.TotalFoods != len(fullSnapshot.Foods) {
		t.Fatalf("expected total food count %d, got %d", len(fullSnapshot.Foods), transportSnapshot.TotalFoods)
	}
	if len(transportSnapshot.MinimapAutonomousCircles) != len(fullSnapshot.AutonomousCircles) {
		t.Fatalf("expected minimap autonomous summaries %d, got %d", len(fullSnapshot.AutonomousCircles), len(transportSnapshot.MinimapAutonomousCircles))
	}
	if len(transportSnapshot.MinimapFoods) != len(fullSnapshot.Foods) {
		t.Fatalf("expected minimap food summaries %d, got %d", len(fullSnapshot.Foods), len(transportSnapshot.MinimapFoods))
	}
}

func TestViewportTransportSnapshotReducesMeasuredPayloadBelowFullSnapshotBaseline(t *testing.T) {
	fullSnapshot := simulation.NewSession().Snapshot()
	transportSnapshot := transport.BuildViewportSnapshot(fullSnapshot)

	fullMeasurement, err := transport.MeasureSnapshotTransport(fullSnapshot, transport.DefaultTickEvery)
	if err != nil {
		t.Fatalf("measure full snapshot transport: %v", err)
	}
	culledMeasurement, err := transport.MeasureSnapshotTransport(transportSnapshot, transport.DefaultTickEvery)
	if err != nil {
		t.Fatalf("measure culled snapshot transport: %v", err)
	}

	if culledMeasurement.PayloadBytes >= fullMeasurement.PayloadBytes {
		t.Fatalf("expected culled payload %d to be smaller than full payload %d", culledMeasurement.PayloadBytes, fullMeasurement.PayloadBytes)
	}
	if culledMeasurement.ApproxBytesPerSecond >= fullMeasurement.ApproxBytesPerSecond {
		t.Fatalf("expected culled bytes/sec %v to be smaller than full bytes/sec %v", culledMeasurement.ApproxBytesPerSecond, fullMeasurement.ApproxBytesPerSecond)
	}
}
