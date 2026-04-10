package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http/httptest"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"

	"github.com/codingzen/living-circles-mrl/src/server/simulation"
)

type SnapshotTransportMeasurement struct {
	PayloadBytes         int
	SnapshotsPerSecond   float64
	ApproxBytesPerSecond float64
}

type ActiveTransportComponentMeasurement struct {
	Full               SnapshotTransportMeasurement
	WithoutPlayer      SnapshotTransportMeasurement
	WithoutAutonomous  SnapshotTransportMeasurement
	WithoutFoods       SnapshotTransportMeasurement
	WithoutOrientation SnapshotTransportMeasurement
	WithoutInteraction SnapshotTransportMeasurement
	DominantComponent  string
}

type ActivePlayerPrecisionCandidateMeasurement struct {
	Base       SnapshotTransportMeasurement
	Candidate  SnapshotTransportMeasurement
	Worthwhile bool
}

type ActiveOrientationUsabilityMeasurement struct {
	Window                    time.Duration
	TotalSnapshots            int
	FreshOrientationSnapshots int
	StaleOrientationSnapshots int
	ExpectedTickEvery         time.Duration
}

type MultiClientTransportMeasurement struct {
	ClientCount                int
	Window                     time.Duration
	PerClientBytes             []int
	PerClientSnapshots         []int
	AggregateBytes             int
	AggregateSnapshots         int
	ApproxAggregateBytesPerSec float64
	ApproxPerClientBytesPerSec float64
	MaxInterSnapshotGap        time.Duration
	ExpectedTickEvery          time.Duration
}

type FanoutScalingMeasurement struct {
	ClientCounts      []int
	Measurements      []MultiClientTransportMeasurement
	Window            time.Duration
	ExpectedTickEvery time.Duration
}

type TwoClientResponsivenessMeasurement struct {
	Window              time.Duration
	ExpectedTickEvery   time.Duration
	OneActiveOneIdle    MultiClientTransportMeasurement
	TwoActive           MultiClientTransportMeasurement
	IdlePathReachable   bool
	ActivePathPreserved bool
}

type TwoActiveTickBroadcastPressureMeasurement struct {
	Window                  time.Duration
	ExpectedTickEvery       time.Duration
	OneActive               MultiClientTransportMeasurement
	TwoActive               MultiClientTransportMeasurement
	GapIncrease             time.Duration
	TimingPressureBounded   bool
	PayloadPressureDominant bool
}

type MultiClientTransportConfig struct {
	ClientCount       int
	Window            time.Duration
	MovingClientCount int
	MovementDirection simulation.Vector
}

func MeasureSnapshotTransport(snapshot any, tickEvery time.Duration) (SnapshotTransportMeasurement, error) {
	if tickEvery <= 0 {
		return SnapshotTransportMeasurement{}, fmt.Errorf("tickEvery must be positive")
	}

	payload, err := json.Marshal(snapshot)
	if err != nil {
		return SnapshotTransportMeasurement{}, err
	}

	snapshotsPerSecond := 1 / tickEvery.Seconds()
	return SnapshotTransportMeasurement{
		PayloadBytes:         len(payload),
		SnapshotsPerSecond:   snapshotsPerSecond,
		ApproxBytesPerSecond: float64(len(payload)) * snapshotsPerSecond,
	}, nil
}

func MeasureActiveTransportComponents(snapshot simulation.Snapshot) (ActiveTransportComponentMeasurement, error) {
	fullSnapshot := BuildViewportSnapshot(snapshot, true)

	full, err := MeasureSnapshotTransport(fullSnapshot, DefaultTickEvery)
	if err != nil {
		return ActiveTransportComponentMeasurement{}, err
	}

	withoutPlayerSnapshot := fullSnapshot
	withoutPlayerSnapshot.Player = nil
	withoutPlayer, err := MeasureSnapshotTransport(withoutPlayerSnapshot, DefaultTickEvery)
	if err != nil {
		return ActiveTransportComponentMeasurement{}, err
	}

	withoutAutonomousSnapshot := fullSnapshot
	withoutAutonomousSnapshot.AutonomousCircles = []simulation.AutonomousCircle{}
	withoutAutonomous, err := MeasureSnapshotTransport(withoutAutonomousSnapshot, DefaultTickEvery)
	if err != nil {
		return ActiveTransportComponentMeasurement{}, err
	}

	withoutFoodsSnapshot := fullSnapshot
	withoutFoodsSnapshot.Foods = []simulation.Food{}
	withoutFoods, err := MeasureSnapshotTransport(withoutFoodsSnapshot, DefaultTickEvery)
	if err != nil {
		return ActiveTransportComponentMeasurement{}, err
	}

	withoutOrientationSnapshot := fullSnapshot
	withoutOrientationSnapshot.OrientationFresh = false
	withoutOrientationSnapshot.MinimapAutonomousCircles = nil
	withoutOrientationSnapshot.MinimapFoods = nil
	withoutOrientation, err := MeasureSnapshotTransport(withoutOrientationSnapshot, DefaultTickEvery)
	if err != nil {
		return ActiveTransportComponentMeasurement{}, err
	}

	withoutInteractionSnapshot := fullSnapshot
	withoutInteractionSnapshot.Interaction = nil
	withoutInteraction, err := MeasureSnapshotTransport(withoutInteractionSnapshot, DefaultTickEvery)
	if err != nil {
		return ActiveTransportComponentMeasurement{}, err
	}

	dominantComponent := dominantActiveComponent(full.PayloadBytes, withoutPlayer.PayloadBytes, withoutAutonomous.PayloadBytes, withoutFoods.PayloadBytes, withoutOrientation.PayloadBytes, withoutInteraction.PayloadBytes)

	return ActiveTransportComponentMeasurement{
		Full:               full,
		WithoutPlayer:      withoutPlayer,
		WithoutAutonomous:  withoutAutonomous,
		WithoutFoods:       withoutFoods,
		WithoutOrientation: withoutOrientation,
		WithoutInteraction: withoutInteraction,
		DominantComponent:  dominantComponent,
	}, nil
}

func MeasureActivePlayerPrecisionCandidate(snapshot simulation.Snapshot) (ActivePlayerPrecisionCandidateMeasurement, error) {
	baseSnapshot := BuildViewportSnapshot(snapshot, true)

	base, err := MeasureSnapshotTransport(baseSnapshot, DefaultTickEvery)
	if err != nil {
		return ActivePlayerPrecisionCandidateMeasurement{}, err
	}

	candidateSnapshot := baseSnapshot
	if candidateSnapshot.Player != nil {
		copy := *candidateSnapshot.Player
		copy.X = coarsePlayerTransportFloat(copy.X)
		copy.Y = coarsePlayerTransportFloat(copy.Y)
		copy.Radius = coarsePlayerTransportFloat(copy.Radius)
		copy.Energy = coarsePlayerTransportFloat(copy.Energy)
		copy.AttachedChildren = coarseRoundedChildren(copy.AttachedChildren)
		candidateSnapshot.Player = &copy
	}

	candidate, err := MeasureSnapshotTransport(candidateSnapshot, DefaultTickEvery)
	if err != nil {
		return ActivePlayerPrecisionCandidateMeasurement{}, err
	}

	return ActivePlayerPrecisionCandidateMeasurement{
		Base:       base,
		Candidate:  candidate,
		Worthwhile: candidate.PayloadBytes < base.PayloadBytes,
	}, nil
}

func MeasureActiveOrientationUsability(session *simulation.Session, window time.Duration, direction simulation.Vector) (ActiveOrientationUsabilityMeasurement, error) {
	if session == nil {
		return ActiveOrientationUsabilityMeasurement{}, fmt.Errorf("session must not be nil")
	}
	if window <= 0 {
		return ActiveOrientationUsabilityMeasurement{}, fmt.Errorf("window must be positive")
	}
	if direction.X == 0 && direction.Y == 0 {
		direction = simulation.Vector{X: 1, Y: 0}
	}

	measurement, err := MeasureMultiClientTransportSnapshots(session, MultiClientTransportConfig{
		ClientCount:       1,
		Window:            window,
		MovingClientCount: 1,
		MovementDirection: direction,
	})
	if err != nil {
		return ActiveOrientationUsabilityMeasurement{}, err
	}

	result := ActiveOrientationUsabilityMeasurement{
		Window:            window,
		ExpectedTickEvery: DefaultTickEvery,
	}
	for _, snapshot := range measurement[0] {
		result.TotalSnapshots++
		if snapshot.OrientationFresh {
			result.FreshOrientationSnapshots++
		} else {
			result.StaleOrientationSnapshots++
		}
	}
	return result, nil
}

func coarsePlayerTransportFloat(value float64) float64 {
	return math.Round(value/4) * 4
}

func coarseRoundedChildren(children []simulation.AttachedChild) []simulation.AttachedChild {
	rounded := make([]simulation.AttachedChild, 0, len(children))
	for _, child := range children {
		copy := child
		copy.X = coarsePlayerTransportFloat(copy.X)
		copy.Y = coarsePlayerTransportFloat(copy.Y)
		copy.Radius = coarsePlayerTransportFloat(copy.Radius)
		rounded = append(rounded, copy)
	}
	return rounded
}

func MeasureTwoClientResponsiveness(sessionFactory func() *simulation.Session, window time.Duration, direction simulation.Vector) (TwoClientResponsivenessMeasurement, error) {
	if sessionFactory == nil {
		return TwoClientResponsivenessMeasurement{}, fmt.Errorf("sessionFactory must not be nil")
	}
	if window <= 0 {
		return TwoClientResponsivenessMeasurement{}, fmt.Errorf("window must be positive")
	}
	if direction.X == 0 && direction.Y == 0 {
		direction = simulation.Vector{X: 1, Y: 0}
	}

	oneActiveOneIdle, err := MeasureMultiClientTransportWithConfig(sessionFactory(), MultiClientTransportConfig{
		ClientCount:       2,
		Window:            window,
		MovingClientCount: 1,
		MovementDirection: direction,
	})
	if err != nil {
		return TwoClientResponsivenessMeasurement{}, err
	}

	twoActive, err := MeasureMultiClientTransportWithConfig(sessionFactory(), MultiClientTransportConfig{
		ClientCount:       2,
		Window:            window,
		MovingClientCount: 2,
		MovementDirection: direction,
	})
	if err != nil {
		return TwoClientResponsivenessMeasurement{}, err
	}

	measurement := TwoClientResponsivenessMeasurement{
		Window:              window,
		ExpectedTickEvery:   DefaultTickEvery,
		OneActiveOneIdle:    oneActiveOneIdle,
		TwoActive:           twoActive,
		IdlePathReachable:   len(oneActiveOneIdle.PerClientSnapshots) >= 2 && oneActiveOneIdle.PerClientSnapshots[1] < oneActiveOneIdle.PerClientSnapshots[0],
		ActivePathPreserved: len(oneActiveOneIdle.PerClientSnapshots) >= 1 && oneActiveOneIdle.PerClientSnapshots[0] >= int(window/DefaultTickEvery),
	}

	return measurement, nil
}

func MeasureTwoActiveTickBroadcastPressure(sessionFactory func() *simulation.Session, window time.Duration, direction simulation.Vector) (TwoActiveTickBroadcastPressureMeasurement, error) {
	if sessionFactory == nil {
		return TwoActiveTickBroadcastPressureMeasurement{}, fmt.Errorf("sessionFactory must not be nil")
	}
	if window <= 0 {
		return TwoActiveTickBroadcastPressureMeasurement{}, fmt.Errorf("window must be positive")
	}
	if direction.X == 0 && direction.Y == 0 {
		direction = simulation.Vector{X: 1, Y: 0}
	}

	oneActive, err := MeasureMultiClientTransportWithConfig(sessionFactory(), MultiClientTransportConfig{
		ClientCount:       1,
		Window:            window,
		MovingClientCount: 1,
		MovementDirection: direction,
	})
	if err != nil {
		return TwoActiveTickBroadcastPressureMeasurement{}, err
	}

	twoActive, err := MeasureMultiClientTransportWithConfig(sessionFactory(), MultiClientTransportConfig{
		ClientCount:       2,
		Window:            window,
		MovingClientCount: 2,
		MovementDirection: direction,
	})
	if err != nil {
		return TwoActiveTickBroadcastPressureMeasurement{}, err
	}

	gapIncrease := twoActive.MaxInterSnapshotGap - oneActive.MaxInterSnapshotGap
	if gapIncrease < 0 {
		gapIncrease = 0
	}

	measurement := TwoActiveTickBroadcastPressureMeasurement{
		Window:                  window,
		ExpectedTickEvery:       DefaultTickEvery,
		OneActive:               oneActive,
		TwoActive:               twoActive,
		GapIncrease:             gapIncrease,
		TimingPressureBounded:   twoActive.MaxInterSnapshotGap <= 2*DefaultTickEvery,
		PayloadPressureDominant: twoActive.AggregateBytes > oneActive.AggregateBytes && twoActive.MaxInterSnapshotGap <= oneActive.MaxInterSnapshotGap+DefaultTickEvery/10,
	}

	return measurement, nil
}

func MeasureMultiClientTransport(session *simulation.Session, clientCount int, window time.Duration) (MultiClientTransportMeasurement, error) {
	return MeasureMultiClientTransportWithConfig(session, MultiClientTransportConfig{
		ClientCount:       clientCount,
		Window:            window,
		MovingClientCount: 0,
	})
}

func MeasureMultiClientTransportSnapshots(session *simulation.Session, config MultiClientTransportConfig) ([][]Snapshot, error) {
	if session == nil {
		return nil, fmt.Errorf("session must not be nil")
	}
	if config.ClientCount <= 0 {
		return nil, fmt.Errorf("clientCount must be positive")
	}
	if config.Window <= 0 {
		return nil, fmt.Errorf("window must be positive")
	}
	if config.MovingClientCount < 0 || config.MovingClientCount > config.ClientCount {
		return nil, fmt.Errorf("movingClientCount must be between 0 and clientCount")
	}
	direction := config.MovementDirection
	if direction.X == 0 && direction.Y == 0 {
		direction = simulation.Vector{X: 1, Y: 0}
	}

	expectedTicks := int(config.Window / DefaultTickEvery)
	if expectedTicks <= 0 {
		return nil, fmt.Errorf("window must span at least one tick")
	}

	server := NewServerWithSession(session)
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"

	type clientSnapshots struct {
		index     int
		snapshots []Snapshot
		err       error
	}

	results := make(chan clientSnapshots, config.ClientCount)
	type connectedClient struct {
		index      int
		connection *websocket.Conn
		initial    Snapshot
	}
	clients := make([]connectedClient, 0, config.ClientCount)

	for clientIndex := range config.ClientCount {
		connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
		if err != nil {
			return nil, err
		}

		var snapshot Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			_ = connection.Close()
			return nil, err
		}

		clients = append(clients, connectedClient{
			index:      clientIndex,
			connection: connection,
			initial:    snapshot,
		})
	}

	go func() {
		done <- server.Run(ctx)
	}()

	var waitGroup sync.WaitGroup
	for _, client := range clients {
		waitGroup.Add(1)
		go func(client connectedClient) {
			defer waitGroup.Done()
			defer client.connection.Close()

			if err := client.connection.SetReadDeadline(time.Now().Add(config.Window + time.Second)); err != nil {
				results <- clientSnapshots{index: client.index, err: err}
				return
			}

			stopSending := make(chan struct{})
			if client.index < config.MovingClientCount {
				if err := sendMovementIntent(client.connection, direction); err != nil {
					results <- clientSnapshots{index: client.index, err: err}
					return
				}
				go keepSendingMovementIntent(stopSending, client.connection, direction)
			}
			defer close(stopSending)

			snapshots := make([]Snapshot, 0, expectedTicks+1)
			snapshots = append(snapshots, client.initial)
			for {
				var snapshot Snapshot
				if err := client.connection.ReadJSON(&snapshot); err != nil {
					if isReadTimeout(err) || isClosedConnection(err) {
						break
					}
					results <- clientSnapshots{index: client.index, err: err}
					return
				}
				snapshots = append(snapshots, snapshot)
			}

			results <- clientSnapshots{index: client.index, snapshots: snapshots}
		}(client)
	}

	waitUntilTick(session, int64(expectedTicks), config.Window+time.Second)

	cancel()
	<-done
	waitGroup.Wait()

	perClient := make([][]Snapshot, config.ClientCount)
	for range config.ClientCount {
		result := <-results
		if result.err != nil {
			return nil, result.err
		}
		perClient[result.index] = append(perClient[result.index], result.snapshots...)
	}

	return perClient, nil
}

func MeasureMultiClientTransportWithConfig(session *simulation.Session, config MultiClientTransportConfig) (MultiClientTransportMeasurement, error) {
	if session == nil {
		return MultiClientTransportMeasurement{}, fmt.Errorf("session must not be nil")
	}
	if config.ClientCount <= 0 {
		return MultiClientTransportMeasurement{}, fmt.Errorf("clientCount must be positive")
	}
	if config.Window <= 0 {
		return MultiClientTransportMeasurement{}, fmt.Errorf("window must be positive")
	}
	if config.MovingClientCount < 0 || config.MovingClientCount > config.ClientCount {
		return MultiClientTransportMeasurement{}, fmt.Errorf("movingClientCount must be between 0 and clientCount")
	}
	direction := config.MovementDirection
	if direction.X == 0 && direction.Y == 0 {
		direction = simulation.Vector{X: 1, Y: 0}
	}

	expectedTicks := int(config.Window / DefaultTickEvery)
	if expectedTicks <= 0 {
		return MultiClientTransportMeasurement{}, fmt.Errorf("window must span at least one tick")
	}

	server := NewServerWithSession(session)
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"

	type clientResult struct {
		index               int
		bytes               int
		snapshots           int
		maxInterSnapshotGap time.Duration
		err                 error
	}

	results := make(chan clientResult, config.ClientCount)
	type connectedClient struct {
		index      int
		connection *websocket.Conn
		bytesRead  int
		snapshots  int
		lastReadAt time.Time
	}
	clients := make([]connectedClient, 0, config.ClientCount)

	for clientIndex := range config.ClientCount {
		connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
		if err != nil {
			return MultiClientTransportMeasurement{}, err
		}

		_, payload, err := connection.ReadMessage()
		if err != nil {
			_ = connection.Close()
			return MultiClientTransportMeasurement{}, err
		}

		clients = append(clients, connectedClient{
			index:      clientIndex,
			connection: connection,
			bytesRead:  len(payload),
			snapshots:  1,
			lastReadAt: time.Now(),
		})
	}

	go func() {
		done <- server.Run(ctx)
	}()

	var waitGroup sync.WaitGroup
	for _, client := range clients {
		waitGroup.Add(1)
		go func(client connectedClient) {
			defer waitGroup.Done()
			defer client.connection.Close()

			if err := client.connection.SetReadDeadline(time.Now().Add(config.Window + time.Second)); err != nil {
				results <- clientResult{index: client.index, err: err}
				return
			}

			stopSending := make(chan struct{})
			if client.index < config.MovingClientCount {
				if err := sendMovementIntent(client.connection, direction); err != nil {
					results <- clientResult{index: client.index, err: err}
					return
				}
				go keepSendingMovementIntent(stopSending, client.connection, direction)
			}
			defer close(stopSending)

			bytesRead := client.bytesRead
			snapshotsRead := client.snapshots
			maxGap := time.Duration(0)
			lastSnapshotAt := client.lastReadAt

			for {
				_, payload, err := client.connection.ReadMessage()
				if err != nil {
					if isReadTimeout(err) || isClosedConnection(err) {
						break
					}
					results <- clientResult{index: client.index, err: err}
					return
				}

				now := time.Now()
				gap := now.Sub(lastSnapshotAt)
				if gap > maxGap {
					maxGap = gap
				}
				lastSnapshotAt = now

				bytesRead += len(payload)
				snapshotsRead++
			}

			results <- clientResult{
				index:               client.index,
				bytes:               bytesRead,
				snapshots:           snapshotsRead,
				maxInterSnapshotGap: maxGap,
			}
		}(client)
	}

	waitUntilTick(session, int64(expectedTicks), config.Window+time.Second)

	cancel()
	<-done
	waitGroup.Wait()

	measurement := MultiClientTransportMeasurement{
		ClientCount:        config.ClientCount,
		Window:             config.Window,
		PerClientBytes:     make([]int, config.ClientCount),
		PerClientSnapshots: make([]int, config.ClientCount),
		ExpectedTickEvery:  DefaultTickEvery,
	}

	for range config.ClientCount {
		result := <-results
		if result.err != nil {
			return MultiClientTransportMeasurement{}, result.err
		}
		measurement.PerClientBytes[result.index] = result.bytes
		measurement.PerClientSnapshots[result.index] = result.snapshots
		measurement.AggregateBytes += result.bytes
		measurement.AggregateSnapshots += result.snapshots
		if result.maxInterSnapshotGap > measurement.MaxInterSnapshotGap {
			measurement.MaxInterSnapshotGap = result.maxInterSnapshotGap
		}
	}

	measurement.ApproxAggregateBytesPerSec = float64(measurement.AggregateBytes) / config.Window.Seconds()
	measurement.ApproxPerClientBytesPerSec = measurement.ApproxAggregateBytesPerSec / float64(config.ClientCount)

	return measurement, nil
}

func sendMovementIntent(connection *websocket.Conn, direction simulation.Vector) error {
	return connection.WriteJSON(movementIntentMessage{
		Type:      "movement_intent",
		Direction: direction,
	})
}

func keepSendingMovementIntent(stop <-chan struct{}, connection *websocket.Conn, direction simulation.Vector) {
	ticker := time.NewTicker(DefaultTickEvery)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			if err := sendMovementIntent(connection, direction); err != nil {
				return
			}
		}
	}
}

func isReadTimeout(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

func isClosedConnection(err error) bool {
	var closeErr *websocket.CloseError
	return websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) ||
		errors.As(err, &closeErr) ||
		errors.Is(err, net.ErrClosed) ||
		errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, io.EOF)
}

func waitUntilTick(session *simulation.Session, targetTick int64, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if session.Snapshot().Tick >= targetTick {
			return
		}
		time.Sleep(time.Millisecond)
	}
}

func MeasureClientCountFanoutScaling(sessionFactory func() *simulation.Session, clientCounts []int, window time.Duration) (FanoutScalingMeasurement, error) {
	if sessionFactory == nil {
		return FanoutScalingMeasurement{}, fmt.Errorf("sessionFactory must not be nil")
	}
	if len(clientCounts) == 0 {
		return FanoutScalingMeasurement{}, fmt.Errorf("clientCounts must not be empty")
	}
	if window <= 0 {
		return FanoutScalingMeasurement{}, fmt.Errorf("window must be positive")
	}

	measurements := make([]MultiClientTransportMeasurement, 0, len(clientCounts))
	for _, clientCount := range clientCounts {
		if clientCount <= 0 {
			return FanoutScalingMeasurement{}, fmt.Errorf("clientCount must be positive")
		}
		measurement, err := MeasureMultiClientTransport(sessionFactory(), clientCount, window)
		if err != nil {
			return FanoutScalingMeasurement{}, err
		}
		measurements = append(measurements, measurement)
	}

	copyCounts := append([]int(nil), clientCounts...)
	return FanoutScalingMeasurement{
		ClientCounts:      copyCounts,
		Measurements:      measurements,
		Window:            window,
		ExpectedTickEvery: DefaultTickEvery,
	}, nil
}

func MeasureActiveClientFanoutScaling(sessionFactory func() *simulation.Session, clientCounts []int, window time.Duration, direction simulation.Vector) (FanoutScalingMeasurement, error) {
	if sessionFactory == nil {
		return FanoutScalingMeasurement{}, fmt.Errorf("sessionFactory must not be nil")
	}
	if len(clientCounts) == 0 {
		return FanoutScalingMeasurement{}, fmt.Errorf("clientCounts must not be empty")
	}
	if window <= 0 {
		return FanoutScalingMeasurement{}, fmt.Errorf("window must be positive")
	}
	if direction.X == 0 && direction.Y == 0 {
		direction = simulation.Vector{X: 1, Y: 0}
	}

	measurements := make([]MultiClientTransportMeasurement, 0, len(clientCounts))
	for _, clientCount := range clientCounts {
		if clientCount <= 0 {
			return FanoutScalingMeasurement{}, fmt.Errorf("clientCount must be positive")
		}
		measurement, err := MeasureMultiClientTransportWithConfig(sessionFactory(), MultiClientTransportConfig{
			ClientCount:       clientCount,
			Window:            window,
			MovingClientCount: clientCount,
			MovementDirection: direction,
		})
		if err != nil {
			return FanoutScalingMeasurement{}, err
		}
		measurements = append(measurements, measurement)
	}

	copyCounts := append([]int(nil), clientCounts...)
	return FanoutScalingMeasurement{
		ClientCounts:      copyCounts,
		Measurements:      measurements,
		Window:            window,
		ExpectedTickEvery: DefaultTickEvery,
	}, nil
}

func dominantActiveComponent(full, withoutPlayer, withoutAutonomous, withoutFoods, withoutOrientation, withoutInteraction int) string {
	type component struct {
		name   string
		saving int
	}
	components := []component{
		{name: "player", saving: full - withoutPlayer},
		{name: "autonomous", saving: full - withoutAutonomous},
		{name: "foods", saving: full - withoutFoods},
		{name: "orientation", saving: full - withoutOrientation},
		{name: "interaction", saving: full - withoutInteraction},
	}
	best := components[0]
	for _, candidate := range components[1:] {
		if candidate.saving > best.saving {
			best = candidate
		}
	}
	return best.name
}
