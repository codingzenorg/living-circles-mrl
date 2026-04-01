package integration_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/codingzen/living-circles-mrl/src/server/simulation"
	"github.com/codingzen/living-circles-mrl/src/server/transport"
)

func TestClientReceivesInitialSnapshotAndInteractionClassification(t *testing.T) {
	server := transport.NewServer()
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	if len(initial.Foods) == 0 {
		t.Fatal("expected initial snapshot to include food items")
	}
	if len(initial.AutonomousCircles) != 1 {
		t.Fatalf("expected one autonomous circle, got %d", len(initial.AutonomousCircles))
	}
	if initial.Player.Shape == "" || initial.AutonomousCircles[0].Shape == "" {
		t.Fatal("expected circle shapes in initial snapshot")
	}

	message := map[string]any{
		"type": "movement_intent",
		"direction": map[string]float64{
			"x": -1,
			"y": 0,
		},
	}
	if err := connection.WriteJSON(message); err != nil {
		t.Fatalf("write movement intent: %v", err)
	}

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read updated snapshot: %v", err)
		}

		if snapshot.Tick == 0 {
			continue
		}

		if snapshot.Interaction != nil {
			expectedKind := "reproduce_candidate"
			if snapshot.Player.Shape == snapshot.AutonomousCircles[0].Shape {
				expectedKind = "fight_candidate"
			}
			if snapshot.Player.X >= initial.Player.X {
				t.Fatalf("expected player to move left toward interaction, before=%v after=%v", initial.Player.X, snapshot.Player.X)
			}
			if snapshot.Interaction.Kind != expectedKind {
				t.Fatalf("expected interaction kind %q, got %q", expectedKind, snapshot.Interaction.Kind)
			}
			return
		}
	}
}
