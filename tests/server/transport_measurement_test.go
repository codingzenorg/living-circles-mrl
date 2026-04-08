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
	transportSnapshot := transport.BuildViewportSnapshot(fullSnapshot, true)

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
	if len(transportSnapshot.MinimapAutonomousCircles) == 0 || len(transportSnapshot.MinimapAutonomousCircles) >= len(fullSnapshot.AutonomousCircles) {
		t.Fatalf("expected compact minimap autonomous summaries to be non-empty and smaller than full world count, compact=%d full=%d", len(transportSnapshot.MinimapAutonomousCircles), len(fullSnapshot.AutonomousCircles))
	}
	if len(transportSnapshot.MinimapFoods) == 0 || len(transportSnapshot.MinimapFoods) >= len(fullSnapshot.Foods) {
		t.Fatalf("expected compact minimap food summaries to be non-empty and smaller than full world count, compact=%d full=%d", len(transportSnapshot.MinimapFoods), len(fullSnapshot.Foods))
	}
	autonomousCount := 0
	for _, cluster := range transportSnapshot.MinimapAutonomousCircles {
		autonomousCount += cluster.Count
	}
	if autonomousCount != len(fullSnapshot.AutonomousCircles) {
		t.Fatalf("expected compact autonomous summary to retain total count %d, got %d", len(fullSnapshot.AutonomousCircles), autonomousCount)
	}
	foodCount := 0
	for _, cluster := range transportSnapshot.MinimapFoods {
		foodCount += cluster.Count
	}
	if foodCount != len(fullSnapshot.Foods) {
		t.Fatalf("expected compact food summary to retain total count %d, got %d", len(fullSnapshot.Foods), foodCount)
	}
}

func TestViewportTransportSnapshotReducesMeasuredPayloadBelowFullSnapshotBaseline(t *testing.T) {
	fullSnapshot := simulation.NewSession().Snapshot()
	transportSnapshot := transport.BuildViewportSnapshot(fullSnapshot, true)

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

func TestViewportTransportSnapshotCanOmitOrientationSummaryOnNonRefreshTicks(t *testing.T) {
	fullSnapshot := simulation.NewSession().Snapshot()
	transportSnapshot := transport.BuildViewportSnapshot(fullSnapshot, false)

	if transportSnapshot.OrientationFresh {
		t.Fatal("expected non-refresh transport snapshot to mark orientation as stale")
	}
	if transportSnapshot.MinimapAutonomousCircles != nil {
		t.Fatalf("expected minimap autonomous summaries to be omitted on non-refresh ticks, got %d", len(transportSnapshot.MinimapAutonomousCircles))
	}
	if transportSnapshot.MinimapFoods != nil {
		t.Fatalf("expected minimap food summaries to be omitted on non-refresh ticks, got %d", len(transportSnapshot.MinimapFoods))
	}
	if transportSnapshot.TotalAutonomousCircles != len(fullSnapshot.AutonomousCircles) {
		t.Fatalf("expected total autonomous count %d, got %d", len(fullSnapshot.AutonomousCircles), transportSnapshot.TotalAutonomousCircles)
	}
	if transportSnapshot.TotalFoods != len(fullSnapshot.Foods) {
		t.Fatalf("expected total food count %d, got %d", len(fullSnapshot.Foods), transportSnapshot.TotalFoods)
	}
}

func TestDualCadenceTransportAverageCostFallsBelowSingleCadenceCulledBaseline(t *testing.T) {
	fullSnapshot := simulation.NewSession().Snapshot()
	singleCadenceMeasurement, err := transport.MeasureSnapshotTransport(transport.BuildViewportSnapshotExactOrientation(fullSnapshot, true), transport.DefaultTickEvery)
	if err != nil {
		t.Fatalf("measure single cadence transport: %v", err)
	}

	totalBytes := 0
	for tick := int64(1); tick <= transport.DefaultOrientationEveryTicks; tick++ {
		snapshot, err := transport.MeasureSnapshotTransport(
			transport.BuildViewportSnapshot(simulation.Snapshot{
				Type:              fullSnapshot.Type,
				Tick:              tick,
				World:             fullSnapshot.World,
				Player:            fullSnapshot.Player,
				AutonomousCircles: fullSnapshot.AutonomousCircles,
				Interaction:       fullSnapshot.Interaction,
				Foods:             fullSnapshot.Foods,
			}, tick%transport.DefaultOrientationEveryTicks == 0),
			transport.DefaultTickEvery,
		)
		if err != nil {
			t.Fatalf("measure dual cadence transport tick %d: %v", tick, err)
		}
		totalBytes += snapshot.PayloadBytes
	}

	averagePayload := float64(totalBytes) / float64(transport.DefaultOrientationEveryTicks)
	averageBytesPerSecond := averagePayload * singleCadenceMeasurement.SnapshotsPerSecond
	if averageBytesPerSecond >= singleCadenceMeasurement.ApproxBytesPerSecond {
		t.Fatalf("expected dual cadence average bytes/sec %v to be below single cadence culled baseline %v", averageBytesPerSecond, singleCadenceMeasurement.ApproxBytesPerSecond)
	}
}

func TestCompactMinimapSummaryReducesOrientationRefreshPayloadBelowExactSummary(t *testing.T) {
	fullSnapshot := simulation.NewSession().Snapshot()

	exactMeasurement, err := transport.MeasureSnapshotTransport(transport.BuildViewportSnapshotExactOrientation(fullSnapshot, true), transport.DefaultTickEvery)
	if err != nil {
		t.Fatalf("measure exact orientation transport: %v", err)
	}
	compactMeasurement, err := transport.MeasureSnapshotTransport(transport.BuildViewportSnapshot(fullSnapshot, true), transport.DefaultTickEvery)
	if err != nil {
		t.Fatalf("measure compact orientation transport: %v", err)
	}

	if compactMeasurement.PayloadBytes >= exactMeasurement.PayloadBytes {
		t.Fatalf("expected compact orientation payload %d to be smaller than exact orientation payload %d", compactMeasurement.PayloadBytes, exactMeasurement.PayloadBytes)
	}
}
