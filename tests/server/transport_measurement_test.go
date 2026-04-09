package server_test

import (
	"testing"
	"time"

	"github.com/codingzen/living-circles-mrl/src/server/simulation"
	"github.com/codingzen/living-circles-mrl/src/server/transport"
)

func expectedFullCadenceSnapshots(window time.Duration) int {
	return int(window/transport.DefaultTickEvery) + 1
}

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

func TestMultiClientTransportMeasurementIsDeterministic(t *testing.T) {
	first, err := transport.MeasureMultiClientTransport(simulation.NewSession(), 4, 300*time.Millisecond)
	if err != nil {
		t.Fatalf("measure first multi-client transport: %v", err)
	}

	second, err := transport.MeasureMultiClientTransport(simulation.NewSession(), 4, 300*time.Millisecond)
	if err != nil {
		t.Fatalf("measure second multi-client transport: %v", err)
	}

	if first.ClientCount != second.ClientCount {
		t.Fatalf("expected same client count, first=%+v second=%+v", first, second)
	}
	if first.AggregateSnapshots != second.AggregateSnapshots {
		t.Fatalf("expected deterministic aggregate snapshot count, first=%+v second=%+v", first, second)
	}
	if first.AggregateBytes != second.AggregateBytes {
		t.Fatalf("expected deterministic aggregate bytes, first=%+v second=%+v", first, second)
	}
	if first.AggregateSnapshots > first.ClientCount && first.MaxInterSnapshotGap <= 0 {
		t.Fatalf("expected positive max inter-snapshot gap once repeated snapshots exist, got %+v", first)
	}
}

func TestMultiClientTransportMeasurementScalesAggregateBytesAbovePerClientBytes(t *testing.T) {
	measurement, err := transport.MeasureMultiClientTransport(simulation.NewSession(), 4, 300*time.Millisecond)
	if err != nil {
		t.Fatalf("measure multi-client transport: %v", err)
	}

	if measurement.AggregateBytes <= 0 {
		t.Fatalf("expected aggregate bytes to be positive, got %+v", measurement)
	}
	if measurement.AggregateSnapshots <= 0 {
		t.Fatalf("expected aggregate snapshots to be positive, got %+v", measurement)
	}
	if measurement.ApproxAggregateBytesPerSec <= measurement.ApproxPerClientBytesPerSec {
		t.Fatalf("expected aggregate bytes/sec %v to exceed per-client bytes/sec %v", measurement.ApproxAggregateBytesPerSec, measurement.ApproxPerClientBytesPerSec)
	}
	for index, bytes := range measurement.PerClientBytes {
		if bytes <= 0 {
			t.Fatalf("expected client %d bytes to be positive, got %+v", index, measurement)
		}
	}
	for index, snapshots := range measurement.PerClientSnapshots {
		if snapshots <= 0 {
			t.Fatalf("expected client %d snapshots to be positive, got %+v", index, measurement)
		}
	}
}

func TestPassiveObserverTransportMeasurementReducesSnapshotCadence(t *testing.T) {
	window := 300 * time.Millisecond
	measurement, err := transport.MeasureMultiClientTransport(simulation.NewSession(), 2, window)
	if err != nil {
		t.Fatalf("measure passive observer pair transport: %v", err)
	}

	for index, snapshots := range measurement.PerClientSnapshots {
		if snapshots >= expectedFullCadenceSnapshots(window) {
			t.Fatalf("expected passive observer %d snapshots to stay below full cadence, got %+v", index, measurement)
		}
	}
	for index, snapshots := range measurement.PerClientSnapshots {
		if snapshots != 1 {
			t.Fatalf("expected calm passive observer %d to keep only the initial snapshot inside the bounded window, got %+v", index, measurement)
		}
	}
}

func TestMovingClientTransportMeasurementKeepsFullCadence(t *testing.T) {
	window := 300 * time.Millisecond
	measurement, err := transport.MeasureMultiClientTransportWithConfig(simulation.NewSession(), transport.MultiClientTransportConfig{
		ClientCount:       1,
		Window:            window,
		MovingClientCount: 1,
		MovementDirection: simulation.Vector{X: 1, Y: 0},
	})
	if err != nil {
		t.Fatalf("measure active client transport: %v", err)
	}

	if measurement.PerClientSnapshots[0] < int(window/transport.DefaultTickEvery) {
		t.Fatalf("expected moving client to stay near full cadence, got %+v", measurement)
	}
}

func TestMultiClientTransportMeasurementReportsTickPressureSignal(t *testing.T) {
	measurement, err := transport.MeasureMultiClientTransport(simulation.NewSession(), 4, 300*time.Millisecond)
	if err != nil {
		t.Fatalf("measure multi-client transport: %v", err)
	}

	if measurement.ExpectedTickEvery != transport.DefaultTickEvery {
		t.Fatalf("expected tick cadence %v, got %+v", transport.DefaultTickEvery, measurement)
	}
	if measurement.AggregateSnapshots == measurement.ClientCount {
		if measurement.MaxInterSnapshotGap != 0 {
			t.Fatalf("expected calm passive observer measurement to report no inter-snapshot gap, got %+v", measurement)
		}
		return
	}
	if measurement.MaxInterSnapshotGap < measurement.ExpectedTickEvery {
		t.Fatalf("expected max inter-snapshot gap %v to be at least one tick interval %v", measurement.MaxInterSnapshotGap, measurement.ExpectedTickEvery)
	}
	if measurement.MaxInterSnapshotGap > 4*measurement.ExpectedTickEvery {
		t.Fatalf("expected max inter-snapshot gap %v to stay within bounded local pressure window", measurement.MaxInterSnapshotGap)
	}
}

func TestMovingMultiClientTransportMeasurementIsDeterministic(t *testing.T) {
	config := transport.MultiClientTransportConfig{
		ClientCount:       4,
		Window:            300 * time.Millisecond,
		MovingClientCount: 1,
		MovementDirection: simulation.Vector{X: 1, Y: 0},
	}

	first, err := transport.MeasureMultiClientTransportWithConfig(simulation.NewSession(), config)
	if err != nil {
		t.Fatalf("measure first moving multi-client transport: %v", err)
	}

	second, err := transport.MeasureMultiClientTransportWithConfig(simulation.NewSession(), config)
	if err != nil {
		t.Fatalf("measure second moving multi-client transport: %v", err)
	}

	if first.AggregateSnapshots != second.AggregateSnapshots {
		t.Fatalf("expected deterministic aggregate snapshot count, first=%+v second=%+v", first, second)
	}
	if first.AggregateBytes != second.AggregateBytes {
		t.Fatalf("expected deterministic aggregate bytes, first=%+v second=%+v", first, second)
	}
}

func TestMovingMultiClientTransportDiffersFromIdleAggregateBytes(t *testing.T) {
	idleMeasurement, err := transport.MeasureMultiClientTransport(simulation.NewSession(), 4, 300*time.Millisecond)
	if err != nil {
		t.Fatalf("measure idle multi-client transport: %v", err)
	}

	movingMeasurement, err := transport.MeasureMultiClientTransportWithConfig(simulation.NewSession(), transport.MultiClientTransportConfig{
		ClientCount:       4,
		Window:            300 * time.Millisecond,
		MovingClientCount: 1,
		MovementDirection: simulation.Vector{X: 1, Y: 0},
	})
	if err != nil {
		t.Fatalf("measure moving multi-client transport: %v", err)
	}

	if movingMeasurement.AggregateBytes == idleMeasurement.AggregateBytes {
		t.Fatalf("expected moving aggregate bytes to differ from idle aggregate bytes, idle=%d moving=%d", idleMeasurement.AggregateBytes, movingMeasurement.AggregateBytes)
	}
	if movingMeasurement.ApproxAggregateBytesPerSec == idleMeasurement.ApproxAggregateBytesPerSec {
		t.Fatalf("expected moving aggregate bytes/sec to differ from idle aggregate bytes/sec, idle=%v moving=%v", idleMeasurement.ApproxAggregateBytesPerSec, movingMeasurement.ApproxAggregateBytesPerSec)
	}
}

func TestMovingMultiClientTransportKeepsBoundedTickPressure(t *testing.T) {
	measurement, err := transport.MeasureMultiClientTransportWithConfig(simulation.NewSession(), transport.MultiClientTransportConfig{
		ClientCount:       4,
		Window:            300 * time.Millisecond,
		MovingClientCount: 1,
		MovementDirection: simulation.Vector{X: 1, Y: 0},
	})
	if err != nil {
		t.Fatalf("measure moving multi-client transport: %v", err)
	}

	if measurement.MaxInterSnapshotGap < measurement.ExpectedTickEvery {
		t.Fatalf("expected moving max inter-snapshot gap %v to be at least one tick interval %v", measurement.MaxInterSnapshotGap, measurement.ExpectedTickEvery)
	}
	if measurement.MaxInterSnapshotGap > 4*measurement.ExpectedTickEvery {
		t.Fatalf("expected moving max inter-snapshot gap %v to stay within bounded local pressure window", measurement.MaxInterSnapshotGap)
	}
}

func TestClientCountFanoutScalingMeasurementIsDeterministic(t *testing.T) {
	first, err := transport.MeasureClientCountFanoutScaling(simulation.NewSession, []int{1, 4, 8}, 300*time.Millisecond)
	if err != nil {
		t.Fatalf("measure first fanout scaling: %v", err)
	}

	second, err := transport.MeasureClientCountFanoutScaling(simulation.NewSession, []int{1, 4, 8}, 300*time.Millisecond)
	if err != nil {
		t.Fatalf("measure second fanout scaling: %v", err)
	}

	if len(first.Measurements) != len(second.Measurements) {
		t.Fatalf("expected same measurement count, first=%+v second=%+v", first, second)
	}
	for index := range first.Measurements {
		if first.Measurements[index].AggregateBytes != second.Measurements[index].AggregateBytes {
			t.Fatalf("expected deterministic aggregate bytes at index %d, first=%+v second=%+v", index, first.Measurements[index], second.Measurements[index])
		}
		if first.Measurements[index].AggregateSnapshots != second.Measurements[index].AggregateSnapshots {
			t.Fatalf("expected deterministic aggregate snapshots at index %d, first=%+v second=%+v", index, first.Measurements[index], second.Measurements[index])
		}
	}
}

func TestClientCountFanoutScalingIncreasesAggregateBytesWithClientCount(t *testing.T) {
	scaling, err := transport.MeasureClientCountFanoutScaling(simulation.NewSession, []int{1, 4, 8}, 300*time.Millisecond)
	if err != nil {
		t.Fatalf("measure fanout scaling: %v", err)
	}

	if len(scaling.Measurements) != 3 {
		t.Fatalf("expected 3 scaling measurements, got %+v", scaling)
	}
	oneClient := scaling.Measurements[0]
	fourClients := scaling.Measurements[1]
	eightClients := scaling.Measurements[2]

	if fourClients.AggregateBytes <= oneClient.AggregateBytes {
		t.Fatalf("expected 4-client aggregate bytes %d to exceed 1-client aggregate bytes %d", fourClients.AggregateBytes, oneClient.AggregateBytes)
	}
	if eightClients.AggregateBytes <= fourClients.AggregateBytes {
		t.Fatalf("expected 8-client aggregate bytes %d to exceed 4-client aggregate bytes %d", eightClients.AggregateBytes, fourClients.AggregateBytes)
	}
	if fourClients.ApproxAggregateBytesPerSec <= oneClient.ApproxAggregateBytesPerSec {
		t.Fatalf("expected 4-client aggregate bytes/sec %v to exceed 1-client aggregate bytes/sec %v", fourClients.ApproxAggregateBytesPerSec, oneClient.ApproxAggregateBytesPerSec)
	}
	if eightClients.ApproxAggregateBytesPerSec <= fourClients.ApproxAggregateBytesPerSec {
		t.Fatalf("expected 8-client aggregate bytes/sec %v to exceed 4-client aggregate bytes/sec %v", eightClients.ApproxAggregateBytesPerSec, fourClients.ApproxAggregateBytesPerSec)
	}
}

func TestClientCountFanoutScalingKeepsPerClientPressureInterpretable(t *testing.T) {
	scaling, err := transport.MeasureClientCountFanoutScaling(simulation.NewSession, []int{1, 4, 8}, 300*time.Millisecond)
	if err != nil {
		t.Fatalf("measure fanout scaling: %v", err)
	}

	for _, measurement := range scaling.Measurements {
		if measurement.ApproxPerClientBytesPerSec <= 0 {
			t.Fatalf("expected per-client bytes/sec to stay positive, got %+v", measurement)
		}
		if measurement.AggregateSnapshots == measurement.ClientCount {
			if measurement.MaxInterSnapshotGap != 0 {
				t.Fatalf("expected calm passive scaling measurement to report no inter-snapshot gap, got %+v", measurement)
			}
			continue
		}
		if measurement.MaxInterSnapshotGap < measurement.ExpectedTickEvery {
			t.Fatalf("expected max gap %v to be at least one tick interval %v", measurement.MaxInterSnapshotGap, measurement.ExpectedTickEvery)
		}
		if measurement.MaxInterSnapshotGap > 4*measurement.ExpectedTickEvery {
			t.Fatalf("expected max gap %v to stay within bounded local pressure window", measurement.MaxInterSnapshotGap)
		}
	}
}

func TestClientCountFanoutScalingDropsBelowPreviousPassiveBaseline(t *testing.T) {
	scaling, err := transport.MeasureClientCountFanoutScaling(simulation.NewSession, []int{1, 4, 8}, 300*time.Millisecond)
	if err != nil {
		t.Fatalf("measure fanout scaling: %v", err)
	}

	fourClients := scaling.Measurements[1]
	if fourClients.AggregateBytes >= 51904 {
		t.Fatalf("expected 4-client passive aggregate bytes %d to drop below prior baseline 51904", fourClients.AggregateBytes)
	}
	if fourClients.AggregateSnapshots >= 16 {
		t.Fatalf("expected 4-client passive aggregate snapshots %d to drop below prior baseline 16", fourClients.AggregateSnapshots)
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

func TestObserverTransportSnapshotKeepsOrientationWhileOmittingLocalDetail(t *testing.T) {
	fullSnapshot := simulation.NewSession().Snapshot()
	transportSnapshot := transport.BuildObserverSnapshot(fullSnapshot, true)

	if transportSnapshot.TransportMode != "observer_orientation_only" {
		t.Fatalf("expected observer transport mode, got %+v", transportSnapshot)
	}
	if transportSnapshot.Player != nil {
		t.Fatalf("expected observer snapshot to omit player detail, got %+v", transportSnapshot.Player)
	}
	if len(transportSnapshot.AutonomousCircles) != 0 {
		t.Fatalf("expected observer snapshot to omit local autonomous detail, got %d", len(transportSnapshot.AutonomousCircles))
	}
	if len(transportSnapshot.Foods) != 0 {
		t.Fatalf("expected observer snapshot to omit local food detail, got %d", len(transportSnapshot.Foods))
	}
	if transportSnapshot.FoodsFresh {
		t.Fatalf("expected observer snapshot foods_fresh to be false, got %+v", transportSnapshot)
	}
	if !transportSnapshot.OrientationFresh {
		t.Fatalf("expected observer snapshot orientation to remain fresh, got %+v", transportSnapshot)
	}
	if len(transportSnapshot.MinimapAutonomousCircles) == 0 {
		t.Fatalf("expected observer snapshot to retain minimap autonomous summaries, got %+v", transportSnapshot)
	}
	if len(transportSnapshot.MinimapFoods) == 0 {
		t.Fatalf("expected observer snapshot to retain minimap food summaries, got %+v", transportSnapshot)
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

func TestViewportTransportSnapshotCanOmitLocalFoodsOnNonRefreshTicks(t *testing.T) {
	fullSnapshot := simulation.NewSession().Snapshot()
	transportSnapshot := transport.BuildViewportSnapshot(fullSnapshot, true)
	transportSnapshot.FoodsFresh = false
	transportSnapshot.Foods = nil

	if transportSnapshot.FoodsFresh {
		t.Fatal("expected non-refresh transport snapshot to mark foods as stale")
	}
	if transportSnapshot.Foods != nil {
		t.Fatalf("expected local foods to be omitted on stale-food ticks, got %d", len(transportSnapshot.Foods))
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

func TestOrientationRefreshPolicySkipsUnchangedSummaryAndFallsBackLater(t *testing.T) {
	fullSnapshot := simulation.NewSession().Snapshot()
	orientationSnapshot := transport.BuildViewportSnapshot(fullSnapshot, true)
	signature := transport.OrientationSummarySignature(orientationSnapshot)

	if !transport.ShouldRefreshOrientation("", 0, 0, signature) {
		t.Fatal("expected empty prior state to force an initial orientation refresh")
	}
	if transport.ShouldRefreshOrientation(signature, 0, transport.DefaultOrientationFallbackTicks-1, signature) {
		t.Fatal("expected unchanged summary to stay stale before the fallback interval")
	}
	if !transport.ShouldRefreshOrientation(signature, 0, transport.DefaultOrientationFallbackTicks, signature) {
		t.Fatal("expected fallback interval to force an orientation refresh")
	}
}

func TestOrientationRefreshPolicyTriggersOnCompactSummaryChange(t *testing.T) {
	baseSnapshot := simulation.NewSession().Snapshot()
	baseTransportSnapshot := transport.BuildViewportSnapshot(baseSnapshot, true)
	baseSignature := transport.OrientationSummarySignature(baseTransportSnapshot)

	changedSnapshot := baseSnapshot
	changedSnapshot.AutonomousCircles = append([]simulation.AutonomousCircle{}, baseSnapshot.AutonomousCircles...)
	changedSnapshot.AutonomousCircles[0].X += transport.DefaultMinimapClusterSize
	changedTransportSnapshot := transport.BuildViewportSnapshot(changedSnapshot, true)
	changedSignature := transport.OrientationSummarySignature(changedTransportSnapshot)

	if changedSignature == baseSignature {
		t.Fatal("expected moved autonomous circle to change the compact orientation summary signature")
	}
	if !transport.ShouldRefreshOrientation(baseSignature, 0, 1, changedSignature) {
		t.Fatal("expected summary change to force an orientation refresh")
	}
}

func TestObserverTransportRefreshPolicySkipsUnchangedSummaryAndFallsBackLater(t *testing.T) {
	fullSnapshot := simulation.NewSession().Snapshot()
	observerSnapshot := transport.BuildObserverSnapshot(fullSnapshot, true)
	signature := transport.ObserverTransportSignature(observerSnapshot)

	if !transport.ShouldRefreshObserverTransport("", 0, 0, signature) {
		t.Fatal("expected empty prior state to force an initial observer refresh")
	}
	if transport.ShouldRefreshObserverTransport(signature, 0, transport.DefaultObserverFallbackTicks-1, signature) {
		t.Fatal("expected unchanged observer summary to stay stale before the fallback interval")
	}
	if !transport.ShouldRefreshObserverTransport(signature, 0, transport.DefaultObserverFallbackTicks, signature) {
		t.Fatal("expected observer fallback interval to force a refresh")
	}
}

func TestObserverTransportRefreshPolicyTriggersOnObserverRelevantChange(t *testing.T) {
	baseSnapshot := simulation.NewSession().Snapshot()
	baseObserverSnapshot := transport.BuildObserverSnapshot(baseSnapshot, true)
	baseSignature := transport.ObserverTransportSignature(baseObserverSnapshot)

	changedSnapshot := baseSnapshot
	changedSnapshot.Interaction = &simulation.InteractionClassification{
		Active:   false,
		Resolved: true,
		Kind:     "fight_resolved",
		SourceID: "circle-2",
		TargetID: "circle-3",
		WinnerID: "circle-2",
		LoserID:  "circle-3",
	}
	changedObserverSnapshot := transport.BuildObserverSnapshot(changedSnapshot, true)
	changedSignature := transport.ObserverTransportSignature(changedObserverSnapshot)

	if changedSignature == baseSignature {
		t.Fatal("expected observer-relevant interaction change to alter observer transport signature")
	}
	if !transport.ShouldRefreshObserverTransport(baseSignature, 0, 1, changedSignature) {
		t.Fatal("expected observer-relevant signature change to force a refresh")
	}
}

func TestEventDrivenOrientationAverageCostFallsBelowFixedDualCadenceBaseline(t *testing.T) {
	fullSnapshot := simulation.NewSession().Snapshot()

	fixedDualCadenceTotalBytes := 0
	for tick := int64(1); tick <= transport.DefaultOrientationFallbackTicks; tick++ {
		includeOrientation := tick%transport.DefaultOrientationEveryTicks == 0
		measurement, err := transport.MeasureSnapshotTransport(
			transport.BuildViewportSnapshot(simulation.Snapshot{
				Type:              fullSnapshot.Type,
				Tick:              tick,
				World:             fullSnapshot.World,
				Player:            fullSnapshot.Player,
				AutonomousCircles: fullSnapshot.AutonomousCircles,
				Interaction:       fullSnapshot.Interaction,
				Foods:             fullSnapshot.Foods,
			}, includeOrientation),
			transport.DefaultTickEvery,
		)
		if err != nil {
			t.Fatalf("measure fixed dual cadence transport tick %d: %v", tick, err)
		}
		fixedDualCadenceTotalBytes += measurement.PayloadBytes
	}

	eventDrivenTotalBytes := 0
	lastSignature := ""
	lastRefreshTick := int64(-1)
	for tick := int64(1); tick <= transport.DefaultOrientationFallbackTicks; tick++ {
		orientationSnapshot := transport.BuildViewportSnapshot(simulation.Snapshot{
			Type:              fullSnapshot.Type,
			Tick:              tick,
			World:             fullSnapshot.World,
			Player:            fullSnapshot.Player,
			AutonomousCircles: fullSnapshot.AutonomousCircles,
			Interaction:       fullSnapshot.Interaction,
			Foods:             fullSnapshot.Foods,
		}, true)
		signature := transport.OrientationSummarySignature(orientationSnapshot)
		includeOrientation := transport.ShouldRefreshOrientation(lastSignature, lastRefreshTick, tick, signature)
		if includeOrientation {
			lastSignature = signature
			lastRefreshTick = tick
		}

		transportSnapshot := orientationSnapshot
		if !includeOrientation {
			transportSnapshot = transport.BuildViewportSnapshot(simulation.Snapshot{
				Type:              fullSnapshot.Type,
				Tick:              tick,
				World:             fullSnapshot.World,
				Player:            fullSnapshot.Player,
				AutonomousCircles: fullSnapshot.AutonomousCircles,
				Interaction:       fullSnapshot.Interaction,
				Foods:             fullSnapshot.Foods,
			}, false)
		}

		measurement, err := transport.MeasureSnapshotTransport(transportSnapshot, transport.DefaultTickEvery)
		if err != nil {
			t.Fatalf("measure event-driven transport tick %d: %v", tick, err)
		}
		eventDrivenTotalBytes += measurement.PayloadBytes
	}

	fixedAverageBytesPerSecond := (float64(fixedDualCadenceTotalBytes) / float64(transport.DefaultOrientationFallbackTicks)) * (1 / transport.DefaultTickEvery.Seconds())
	eventDrivenAverageBytesPerSecond := (float64(eventDrivenTotalBytes) / float64(transport.DefaultOrientationFallbackTicks)) * (1 / transport.DefaultTickEvery.Seconds())
	if eventDrivenAverageBytesPerSecond >= fixedAverageBytesPerSecond {
		t.Fatalf("expected event-driven average bytes/sec %v to be below fixed dual-cadence baseline %v", eventDrivenAverageBytesPerSecond, fixedAverageBytesPerSecond)
	}
}

func TestLocalFoodRefreshPolicySkipsUnchangedFoodsAndFallsBackLater(t *testing.T) {
	fullSnapshot := simulation.NewSession().Snapshot()
	transportSnapshot := transport.BuildViewportSnapshot(fullSnapshot, true)
	signature := transport.LocalFoodSignature(transportSnapshot)

	if !transport.ShouldRefreshLocalFoods("", 0, 0, signature) {
		t.Fatal("expected empty prior state to force an initial local-food refresh")
	}
	if transport.ShouldRefreshLocalFoods(signature, 0, transport.DefaultLocalFoodFallbackTicks-1, signature) {
		t.Fatal("expected unchanged local food detail to stay stale before the fallback interval")
	}
	if !transport.ShouldRefreshLocalFoods(signature, 0, transport.DefaultLocalFoodFallbackTicks, signature) {
		t.Fatal("expected fallback interval to force a local-food refresh")
	}
}

func TestLocalFoodRefreshPolicyTriggersOnVisibleFoodChange(t *testing.T) {
	baseSnapshot := simulation.NewSession().Snapshot()
	baseTransportSnapshot := transport.BuildViewportSnapshot(baseSnapshot, true)
	baseSignature := transport.LocalFoodSignature(baseTransportSnapshot)
	if len(baseTransportSnapshot.Foods) == 0 {
		t.Fatal("expected base transport snapshot to include at least one visible food")
	}

	changedSnapshot := baseSnapshot
	changedSnapshot.Foods = append([]simulation.Food{}, baseSnapshot.Foods...)
	visibleFoodID := baseTransportSnapshot.Foods[0].ID
	for index, food := range changedSnapshot.Foods {
		if food.ID == visibleFoodID {
			changedSnapshot.Foods[index].X += 120
			break
		}
	}
	changedTransportSnapshot := transport.BuildViewportSnapshot(changedSnapshot, true)
	changedSignature := transport.LocalFoodSignature(changedTransportSnapshot)

	if changedSignature == baseSignature {
		t.Fatal("expected moved visible food to change the local food signature")
	}
	if !transport.ShouldRefreshLocalFoods(baseSignature, 0, 1, changedSignature) {
		t.Fatal("expected visible local food change to force a refresh")
	}
}

func TestEventDrivenLocalFoodAverageCostFallsBelowEventDrivenOrientationBaseline(t *testing.T) {
	fullSnapshot := simulation.NewSession().Snapshot()

	baselineTotalBytes := 0
	lastOrientationSignature := ""
	lastOrientationTick := int64(-1)
	for tick := int64(1); tick <= transport.DefaultLocalFoodFallbackTicks; tick++ {
		orientationSnapshot := transport.BuildViewportSnapshot(simulation.Snapshot{
			Type:              fullSnapshot.Type,
			Tick:              tick,
			World:             fullSnapshot.World,
			Player:            fullSnapshot.Player,
			AutonomousCircles: fullSnapshot.AutonomousCircles,
			Interaction:       fullSnapshot.Interaction,
			Foods:             fullSnapshot.Foods,
		}, true)
		orientationSignature := transport.OrientationSummarySignature(orientationSnapshot)
		includeOrientation := transport.ShouldRefreshOrientation(lastOrientationSignature, lastOrientationTick, tick, orientationSignature)
		if includeOrientation {
			lastOrientationSignature = orientationSignature
			lastOrientationTick = tick
		} else {
			orientationSnapshot.OrientationFresh = false
			orientationSnapshot.MinimapAutonomousCircles = nil
			orientationSnapshot.MinimapFoods = nil
		}

		measurement, err := transport.MeasureSnapshotTransport(orientationSnapshot, transport.DefaultTickEvery)
		if err != nil {
			t.Fatalf("measure orientation-only event-driven baseline tick %d: %v", tick, err)
		}
		baselineTotalBytes += measurement.PayloadBytes
	}

	eventDrivenFoodTotalBytes := 0
	lastOrientationSignature = ""
	lastOrientationTick = int64(-1)
	lastFoodSignature := ""
	lastFoodTick := int64(-1)
	for tick := int64(1); tick <= transport.DefaultLocalFoodFallbackTicks; tick++ {
		transportSnapshot := transport.BuildViewportSnapshot(simulation.Snapshot{
			Type:              fullSnapshot.Type,
			Tick:              tick,
			World:             fullSnapshot.World,
			Player:            fullSnapshot.Player,
			AutonomousCircles: fullSnapshot.AutonomousCircles,
			Interaction:       fullSnapshot.Interaction,
			Foods:             fullSnapshot.Foods,
		}, true)

		orientationSignature := transport.OrientationSummarySignature(transportSnapshot)
		includeOrientation := transport.ShouldRefreshOrientation(lastOrientationSignature, lastOrientationTick, tick, orientationSignature)
		if includeOrientation {
			lastOrientationSignature = orientationSignature
			lastOrientationTick = tick
		} else {
			transportSnapshot.OrientationFresh = false
			transportSnapshot.MinimapAutonomousCircles = nil
			transportSnapshot.MinimapFoods = nil
		}

		foodSignature := transport.LocalFoodSignature(transportSnapshot)
		includeFoods := transport.ShouldRefreshLocalFoods(lastFoodSignature, lastFoodTick, tick, foodSignature)
		if includeFoods {
			lastFoodSignature = foodSignature
			lastFoodTick = tick
		} else {
			transportSnapshot.FoodsFresh = false
			transportSnapshot.Foods = nil
		}

		measurement, err := transport.MeasureSnapshotTransport(transportSnapshot, transport.DefaultTickEvery)
		if err != nil {
			t.Fatalf("measure event-driven local food transport tick %d: %v", tick, err)
		}
		eventDrivenFoodTotalBytes += measurement.PayloadBytes
	}

	baselineAverageBytesPerSecond := (float64(baselineTotalBytes) / float64(transport.DefaultLocalFoodFallbackTicks)) * (1 / transport.DefaultTickEvery.Seconds())
	eventDrivenFoodAverageBytesPerSecond := (float64(eventDrivenFoodTotalBytes) / float64(transport.DefaultLocalFoodFallbackTicks)) * (1 / transport.DefaultTickEvery.Seconds())
	if eventDrivenFoodAverageBytesPerSecond >= baselineAverageBytesPerSecond {
		t.Fatalf("expected event-driven local-food bytes/sec %v to be below event-driven orientation baseline %v", eventDrivenFoodAverageBytesPerSecond, baselineAverageBytesPerSecond)
	}
}

func TestPassiveFanoutScalingDropsBelowObserverOrientationBaseline(t *testing.T) {
	scaling, err := transport.MeasureClientCountFanoutScaling(simulation.NewSession, []int{1, 4, 8}, 300*time.Millisecond)
	if err != nil {
		t.Fatalf("measure fanout scaling: %v", err)
	}

	fourClients := scaling.Measurements[1]
	if fourClients.AggregateBytes >= 22988 {
		t.Fatalf("expected 4-client passive aggregate bytes %d to drop below prior observer-orientation baseline 22988", fourClients.AggregateBytes)
	}
	if fourClients.AggregateSnapshots >= 8 {
		t.Fatalf("expected 4-client passive aggregate snapshots %d to drop below prior observer-orientation baseline 8", fourClients.AggregateSnapshots)
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

func TestReducedLocalPrecisionLowersCompactSummaryPayload(t *testing.T) {
	fullSnapshot := simulation.NewSession().Snapshot()

	fullPrecisionMeasurement, err := transport.MeasureSnapshotTransport(transport.BuildViewportSnapshotCompactFullPrecision(fullSnapshot, true), transport.DefaultTickEvery)
	if err != nil {
		t.Fatalf("measure full-precision compact transport: %v", err)
	}
	reducedPrecisionMeasurement, err := transport.MeasureSnapshotTransport(transport.BuildViewportSnapshot(fullSnapshot, true), transport.DefaultTickEvery)
	if err != nil {
		t.Fatalf("measure reduced-precision compact transport: %v", err)
	}

	if reducedPrecisionMeasurement.PayloadBytes >= fullPrecisionMeasurement.PayloadBytes {
		t.Fatalf("expected reduced-precision payload %d to be smaller than full-precision payload %d", reducedPrecisionMeasurement.PayloadBytes, fullPrecisionMeasurement.PayloadBytes)
	}
}
