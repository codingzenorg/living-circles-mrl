package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/codingzen/living-circles-mrl/src/server/simulation"
)

type SnapshotTransportMeasurement struct {
	PayloadBytes         int
	SnapshotsPerSecond   float64
	ApproxBytesPerSecond float64
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

func MeasureMultiClientTransport(session *simulation.Session, clientCount int, window time.Duration) (MultiClientTransportMeasurement, error) {
	return MeasureMultiClientTransportWithConfig(session, MultiClientTransportConfig{
		ClientCount:       clientCount,
		Window:            window,
		MovingClientCount: 0,
	})
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
	go func() {
		done <- server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"

	type clientResult struct {
		index               int
		bytes               int
		snapshots           int
		maxInterSnapshotGap time.Duration
		err                 error
	}

	results := make(chan clientResult, config.ClientCount)
	var waitGroup sync.WaitGroup

	for clientIndex := range config.ClientCount {
		waitGroup.Add(1)
		go func(index int) {
			defer waitGroup.Done()

			connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
			if err != nil {
				results <- clientResult{index: index, err: err}
				return
			}
			defer connection.Close()

			if err := connection.SetReadDeadline(time.Now().Add(config.Window + time.Second)); err != nil {
				results <- clientResult{index: index, err: err}
				return
			}

			stopSending := make(chan struct{})
			if index < config.MovingClientCount {
				if err := sendMovementIntent(connection, direction); err != nil {
					results <- clientResult{index: index, err: err}
					return
				}
				go keepSendingMovementIntent(stopSending, connection, direction)
			}
			defer close(stopSending)

			bytesRead := 0
			snapshotsRead := 0
			maxGap := time.Duration(0)
			var lastSnapshotAt time.Time
			expectedSnapshots := expectedPassiveObserverSnapshots(expectedTicks, config.ClientCount > 1)
			if index < config.MovingClientCount {
				expectedSnapshots = expectedTicks + 1
			}

			for range expectedSnapshots {
				_, payload, err := connection.ReadMessage()
				if err != nil {
					results <- clientResult{index: index, err: err}
					return
				}

				now := time.Now()
				if !lastSnapshotAt.IsZero() {
					gap := now.Sub(lastSnapshotAt)
					if gap > maxGap {
						maxGap = gap
					}
				}
				lastSnapshotAt = now

				bytesRead += len(payload)
				snapshotsRead++
			}

			results <- clientResult{
				index:               index,
				bytes:               bytesRead,
				snapshots:           snapshotsRead,
				maxInterSnapshotGap: maxGap,
			}
		}(clientIndex)
	}

	waitGroup.Wait()
	cancel()
	<-done

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

func expectedPassiveObserverSnapshots(expectedTicks int, reducePassiveCadence bool) int {
	if !reducePassiveCadence {
		return expectedTicks + 1
	}

	snapshots := 1
	for tick := 1; tick <= expectedTicks; tick++ {
		if tick%int(DefaultPassiveObserverCadenceTicks) == 0 {
			snapshots++
		}
	}

	return snapshots
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
