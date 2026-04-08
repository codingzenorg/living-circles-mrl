package transport

import (
	"encoding/json"
	"fmt"
	"time"
)

type SnapshotTransportMeasurement struct {
	PayloadBytes         int
	SnapshotsPerSecond   float64
	ApproxBytesPerSecond float64
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
